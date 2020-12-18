all: build

init:           ## Installs tools to $GOPATH/bin (which is expected to be in $PATH).
	curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	go install ./vendor/github.com/jstemmer/go-junit-report

build:
	go install -v ./...
	go test -i -v ./...
	go test -c -v ./inventory
	go test -c -v ./management
	go test -c -v ./server

dev-test:						## Run test on dev env. Use `PMM_KUBECONFIG=/path/to/kubeconfig.yaml make dev-test` to run tests for DBaaS.
	go test -count=1 -p 1 -v ./... -pmm.server-insecure-tls

run:
	go test -count=1 -p 1 -v ./... 2>&1 | tee pmm-api-tests-output.txt
	cat pmm-api-tests-output.txt | go-junit-report > pmm-api-tests-junit-report.xml

run-race:
	go test -count=1 -p 1 -v -race ./... 2>&1 | tee pmm-api-tests-output.txt
	cat pmm-api-tests-output.txt | go-junit-report > pmm-api-tests-junit-report.xml

FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:                         ## Format source code.
	gofmt -w -s $(FILES)
	goimports -local github.com/Percona-Lab/pmm-api-tests -l -w $(FILES)

clean:
	rm -f ./pmm-api-tests-output.txt
	rm -f ./pmm-api-tests-junit-report.xml

check-all:                      ## Run golang ci linter to check new changes from master.
	golangci-lint run -c=.golangci.yml --new-from-rev=master

ci-reviewdog:                   ## Runs reviewdog checks.
	golangci-lint run -c=.golangci-required.yml --out-format=line-number | bin/reviewdog -f=golangci-lint -level=error -reporter=github-pr-check
	golangci-lint run -c=.golangci.yml --out-format=line-number | bin/reviewdog -f=golangci-lint -level=error -reporter=github-pr-review
