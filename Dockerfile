FROM golang:1.15

RUN mkdir -p $GOPATH/src/github.com/Percona-Lab/pmm-api-tests

WORKDIR $GOPATH/src/github.com/Percona-Lab/pmm-api-tests/
COPY . $GOPATH/src/github.com/Percona-Lab/pmm-api-tests/

CMD make init run-race
