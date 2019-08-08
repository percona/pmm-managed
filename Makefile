build:
	go install -v ./...
	go test -i -v ./...
	go test -c -v ./inventory
	go test -c -v ./management
	go test -c -v ./server

run:
	go test -count=1 -v ./...

run-race:
	go test -count=1 -v -race ./...

FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:                         ## Format source code.
	gofmt -w -s $(FILES)
	goimports -local github.com/Percona-Lab/pmm-api-tests -l -w $(FILES)
