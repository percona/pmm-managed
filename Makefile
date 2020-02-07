help:                           ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

# `cut` is used to remove first `v` from `git describe` output
PMM_RELEASE_PATH ?= bin
PMM_RELEASE_VERSION ?= $(shell git describe --always --dirty | cut -b2-)
PMM_RELEASE_TIMESTAMP ?= $(shell date '+%s')
PMM_RELEASE_FULLCOMMIT ?= $(shell git rev-parse HEAD)
PMM_RELEASE_BRANCH ?= $(shell git describe --always --contains --all)

LD_FLAGS = -ldflags " \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.ProjectName=pmm-managed' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Version=$(PMM_RELEASE_VERSION)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.PMMVersion=$(PMM_RELEASE_VERSION)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Timestamp=$(PMM_RELEASE_TIMESTAMP)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.FullCommit=$(PMM_RELEASE_FULLCOMMIT)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Branch=$(PMM_RELEASE_BRANCH)' \
			"

release:                        ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed

init:                           ## Installs tools to $GOPATH/bin (which is expected to be in $PATH).
	curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

	# install the same version as a version of Prometheus
	-promtool --version
	mkdir -p /tmp/promtool
	-cd /tmp/promtool && curl -L -O -J https://github.com/prometheus/prometheus/releases/download/v2.12.0/prometheus-2.12.0.$(shell go env GOOS)-amd64.tar.gz
	cd /tmp/promtool/ && tar xvf prometheus-2.12.0.$(shell go env GOOS)-amd64.tar.gz --strip-components 1
	mv /tmp/promtool/promtool $(GOPATH)/bin
	which -a promtool
	promtool --version

	go install ./vendor/github.com/BurntSushi/go-sumtype \
				./vendor/github.com/vektra/mockery/cmd/mockery \
				./vendor/golang.org/x/tools/cmd/goimports \
				./vendor/gopkg.in/reform.v1/reform

	go test -i ./...
	go test -race -i ./...

gen:                            ## Generate files.
	rm -f models/*_reform.go
	find . -name mock_*.go -delete
	go generate ./...
	make format

install:                        ## Install pmm-managed binary.
	go install $(LD_FLAGS) ./...

install-race:                   ## Install pmm-managed binary with race detector.
	go install $(LD_FLAGS) -race ./...

TEST_PACKAGES ?= ./...
TEST_FLAGS ?= -timeout=30s
TEST_RUN_UPDATE ?= 0

test:                           ## Run tests.
	go test $(TEST_FLAGS) -p 1 $(TEST_PACKAGES)

test-race:                      ## Run tests with race detector.
	go test $(TEST_FLAGS) -p 1 -race $(TEST_PACKAGES)

test-cover:                     ## Run tests and collect per-package coverage information.
	go test $(TEST_FLAGS) -p 1 -coverprofile=cover.out -covermode=count $(TEST_PACKAGES)

test-crosscover:                ## Run tests and collect cross-package coverage information.
	go test $(TEST_FLAGS) -p 1 -coverprofile=crosscover.out -covermode=count -coverpkg=./... $(TEST_PACKAGES)

test-race-crosscover:           ## Run tests with race detector and collect cross-package coverage information.
	go test $(TEST_FLAGS) -p 1 -race -coverprofile=race-crosscover.out -covermode=atomic -coverpkg=./... $(TEST_PACKAGES)

fuzz-grafana:                   ## Run fuzzer for services/grafana package.
	# go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build
	mkdir -p services/grafana/fuzzdata/corpus
	cd services/grafana && go-fuzz-build
	cd services/grafana && go-fuzz -workdir=fuzzdata

check:                          ## Run required checkers and linters.
	go run .github/check-license.go
	go-sumtype ./vendor/... ./...

FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:                         ## Format source code.
	gofmt -w -s $(FILES)
	goimports -local github.com/percona/pmm-managed -l -w $(FILES)

RUN_FLAGS = --debug \
			--prometheus-config=testdata/prometheus/prometheus.yml \
			--postgres-name=pmm-managed-dev \
			--supervisord-config-dir=testdata/supervisord.d \
			--alert-manager-rules-file=/tmp/pmm.rules.yml

run: install _run               ## Run pmm-managed.

run-race: install-race _run     ## Run pmm-managed with race detector.

run-race-cover: install-race    ## Run pmm-managed with race detector and collect coverage information.
	go test -coverpkg="github.com/percona/pmm-managed/..." \
			-tags maincover \
			$(LD_FLAGS) \
			-race -c -o bin/pmm-managed.test
	bin/pmm-managed.test -test.coverprofile=cover.out -test.run=TestMainCover $(RUN_FLAGS)

_run:
	pmm-managed $(RUN_FLAGS)

devcontainer:                   ## Run TARGET in devcontainer.
	docker exec pmm-managed-server env \
		TEST_FLAGS='$(TEST_FLAGS)' \
		TEST_PACKAGES='$(TEST_PACKAGES)' \
		TEST_RUN_UPDATE=$(TEST_RUN_UPDATE) \
		make -C /root/go/src/github.com/percona/pmm-managed $(TARGET)

env-up:                         ## Start development environment.
	docker-compose up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans

env-down:                       ## Stop development environment.
	docker-compose down --volumes --remove-orphans

env-psql:                       ## Open psql shell.
	env PGPASSWORD=pmm-managed psql -h 127.0.0.1 -p 5432 -U pmm-managed pmm-managed-dev

clean:                          ## Removes generated artifacts.
	rm -Rf ./bin
