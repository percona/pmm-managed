# This Dockerfile used only for the API tests.

FROM golang:1.18

RUN mkdir -p $GOPATH/src/github.com/percona/pmm-managed/api-tests
RUN mkdir -p $GOPATH/src/github.com/percona/pmm-managed/api-tests-unmocked

WORKDIR $GOPATH/src/github.com/percona/pmm-managed
COPY api-tests/ $GOPATH/src/github.com/percona/pmm-managed/api-tests/
COPY api-tests-unmocked/ $GOPATH/src/github.com/percona/pmm-managed/api-tests-unmocked/
COPY go.mod go.sum $GOPATH/src/github.com/percona/pmm-managed/

CMD make init run-race
