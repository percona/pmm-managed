# pmm-api

[![Build Status](https://travis-ci.org/Percona-Lab/pmm-api.svg?branch=master)](https://travis-ci.org/Percona-Lab/pmm-api)

pmm-api prototype for PMM 2.0.

## Local setup

Generate TLS certificate for `nginx` for local testing:
```
brew install mkcert
mkcert -install
make cert
```

Install `prototool` and fill `vendor/`:
```
make init
```

Generate files:
```
make gen
```

Serve API documentation with `nginx`:
```
make serve
```
