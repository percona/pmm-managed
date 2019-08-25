all: build

init:
	go install ./vendor/github.com/jstemmer/go-junit-report

build:
	go install -v ./...
	go test -i -v ./...
	go test -c -v ./inventory
	go test -c -v ./management
	go test -c -v ./server

run:
	go test -count=1 -v ./... 2>&1 | tee pmm-api-tests-output.txt
	cat pmm-api-tests-output.txt | go-junit-report > pmm-api-tests-junit-report.xml

run-race:
	go test -count=1 -v -race ./... 2>&1 | tee pmm-api-tests-output.txt
	cat pmm-api-tests-output.txt | go-junit-report > pmm-api-tests-junit-report.xml

FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:                         ## Format source code.
	gofmt -w -s $(FILES)
	goimports -local github.com/Percona-Lab/pmm-api-tests -l -w $(FILES)

clean:
	rm -f ./pmm-api-tests-output.txt
	rm -f ./pmm-api-tests-junit-report.xml
