module github.com/percona/pmm-managed/api-tests

go 1.16

// Use for local development, but do not commit:
// replace github.com/percona/pmm => ../../pmm

// Update with:
// go get -v github.com/percona/pmm@PMM-2.0

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/brianvoe/gofakeit/v6 v6.2.2
	github.com/davecgh/go-spew v1.1.1
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/percona-platform/saas v0.0.0-20210122115803-1b32ca1828e1
	github.com/percona/pmm v2.15.2-0.20210408121032-ab1dbcacb092+incompatible
	github.com/prometheus/client_golang v1.9.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/mongo-driver v1.3.5 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
