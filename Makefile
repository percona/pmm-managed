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

env-up:                         ## Start development environment.
	docker-compose pull
	docker-compose up --detach --force-recreate --renew-anon-volumes --remove-orphans
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server .devcontainer/setup.py
	docker exec pmm-managed-server env

env-down:                       ## Stop development environment.
	docker-compose down --volumes --remove-orphans

devcontainer:                   ## Run TARGET in devcontainer.
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed \
		--env TEST_FLAGS='$(TEST_FLAGS)' \
		--env TEST_PACKAGES='$(TEST_PACKAGES)' \
		--env TEST_RUN_UPDATE=$(TEST_RUN_UPDATE) \
		pmm-managed-server make $(TARGET)

release:                        ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed

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

psql:                           ## Open psql shell.
	env PGPASSWORD=pmm-managed psql -U pmm-managed pmm-managed-dev
