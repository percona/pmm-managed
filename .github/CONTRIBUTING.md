# Contributing notes

## Local setup

Run `make init` to install dependencies.

You should also have `mysqld_exporter` and `rds_exporter` binaries somewhere in you `$PATH`.
One way to get them is to install them using `go get`:
```sh
go get -u github.com/percona/mysqld_exporter
go get -u github.com/percona/rds_exporter
cd $GOPATH/src/github.com/percona
git clone https://github.com/percona/postgres_exporter.git
cd postgres_exporter && make build && cp postgres_exporter $GOPATH/bin/
```

You have to use Docker Compose to run most of the tests.

```sh
make up
```

```sh
make
```

Start pmm-managed with

```sh
make run
```

Swagger UI will be available on http://127.0.0.1:7772/swagger/.

## Vendoring

We use [dep](https://github.com/golang/dep) to vendor dependencies.
