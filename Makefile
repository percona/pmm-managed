# Host Makefile.

include Makefile.include

env-up: env-compose-up env-devcontainer     ## Start devcontainer.

env-compose-up:
	docker-compose pull
	docker-compose up --detach --force-recreate --renew-anon-volumes --remove-orphans -e CI_REPO_OWNER -e CI_REPO_NAME -e CI_PULL_REQUEST -e CI_PULL_REQUEST -e CI_COMMIT -e CI_BRANCH -e secrets.GITHUB_TOKEN

env-devcontainer:
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server .devcontainer/setup.py

env-down:                                   ## Stop devcontainer.
	docker-compose down --volumes --remove-orphans

TARGET ?= _bash

env:                                        ## Run `make TARGET` in devcontainer (`make env TARGET=help`); TARGET defaults to bash.
	docker exec -it --workdir=/root/go/src/github.com/percona/pmm-managed pmm-managed-server make $(TARGET)
