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
