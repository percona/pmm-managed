FROM golang:latest

RUN mkdir -p $GOPATH/src/github.com/Percona-Lab/pmm-api-tests

WORKDIR $GOPATH/src/github.com/Percona-Lab/pmm-api-tests/
COPY . $GOPATH/src/github.com/Percona-Lab/pmm-api-tests/

CMD make run
