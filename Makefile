# Host Makefile.

include Makefile.include

release:                        ## Build pmm-managed release binary.
	env CGO_ENABLED=0 go build -v $(PMM_LD_FLAGS) -o $(PMM_RELEASE_PATH)/pmm-managed

env-up:                         ## Start devcontainer.
	docker-compose pull
	docker-compose up --detach --force-recreate --renew-anon-volumes --remove-orphans
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server .devcontainer/setup.py

env-down:                       ## Stop devcontainer.
	docker-compose down --volumes --remove-orphans

TARGET ?= _bash

devcontainer:                   ## Run `make TARGET` in devcontainer (`make devcontainer TARGET=test`); TARGET defaults to bash.
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server make $(TARGET)
