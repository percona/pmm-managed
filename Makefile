# Host Makefile.

# `cut` is used to remove first `v` from `git describe` output
PMM_RELEASE_PATH ?= bin
PMM_RELEASE_VERSION ?= $(shell git describe --always --dirty | cut -b2-)
PMM_RELEASE_TIMESTAMP ?= $(shell date '+%s')
PMM_RELEASE_FULLCOMMIT ?= $(shell git rev-parse HEAD)
PMM_RELEASE_BRANCH ?= $(shell git describe --always --contains --all)

PMM_LD_FLAGS = -ldflags " \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.ProjectName=pmm-managed' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Version=$(PMM_RELEASE_VERSION)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.PMMVersion=$(PMM_RELEASE_VERSION)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Timestamp=$(PMM_RELEASE_TIMESTAMP)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.FullCommit=$(PMM_RELEASE_FULLCOMMIT)' \
			-X 'github.com/percona/pmm-managed/vendor/github.com/percona/pmm/version.Branch=$(PMM_RELEASE_BRANCH)' \
			"

help:                           ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

release:                        ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed

init:                           ## Installs tools to $GOPATH/bin (which is expected to be in $PATH).
	curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

	go install ./vendor/github.com/BurntSushi/go-sumtype \
				./vendor/github.com/vektra/mockery/cmd/mockery \
				./vendor/golang.org/x/tools/cmd/goimports \
				./vendor/gopkg.in/reform.v1/reform

	go get -u github.com/prometheus/prometheus/cmd/promtool

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
			--supervisord-config-dir=testdata/supervisord.d

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

env-up:                         ## Start devcontainer.
	docker-compose pull
	docker-compose up --detach --force-recreate --renew-anon-volumes --remove-orphans
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server .devcontainer/setup.py

env-down:                       ## Stop devcontainer.
	docker-compose down --volumes --remove-orphans

TARGET ?= _bash

devcontainer:                   ## Run `make TARGET` in devcontainer (`make devcontainer TARGET=test`); TARGET defaults to bash.
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server make $(TARGET)

release:                        ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(PMM_LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed
