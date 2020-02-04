# Host Makefile.

include Makefile.include

release:                                    ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(PMM_LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed

env-up: env-compose-up env-devcontainer     ## Start devcontainer.

env-compose-up:
	docker-compose pull
	docker-compose up --detach --force-recreate --renew-anon-volumes --remove-orphans

env-devcontainer:
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server .devcontainer/setup.py

env-down:                                   ## Stop devcontainer.
	docker-compose down --volumes --remove-orphans

TARGET ?= _bash

devcontainer:                               ## Run `make TARGET` in devcontainer (`make devcontainer TARGET=help`); TARGET defaults to bash.
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server make $(TARGET)
