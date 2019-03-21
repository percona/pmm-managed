# pmm-api-tests

API tests for PMM 2.x

# Setup Instructions

Make sure you have Go 1.12.x installed on your systems, execute the following steps
to setup API-tests in your local systems.

1. Fetch the Repo: `go get -u -v github.com/Percona-Lab/pmm-api-tests`
2. Navigate to the tests root folder:  `cd ~/go/src/github.com/Percona-Lab/pmm-api-tests`
3. `make`

# Usage

Once the binaries for the tests have been generated, look at the usage of the tests using:
```
./inventory.test -h
```

It should provide a list of options that are available for execution.

Run the tests using the following command:

```
./inventory.test -pmm.server-url **pmm-server-url** -test.v
```

where `pmm-server-url` should be pointing to pmm-server.
