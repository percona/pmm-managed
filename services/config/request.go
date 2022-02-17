package config

import (
	"github.com/percona/pmm-managed/utils/saasreq"
	"net/http"
)

func (c *Config) configureSaasReqEnrichment() {
	if c.ExtraHeaders.Enabled {
		saasreq.RequestEnricher = func(method string, endpoint string, req *http.Request) {
			for _, each := range c.ExtraHeaders.Endpoints {
				if each.Method == method && each.Endpoint == endpoint {
					next := req.Header
					for name, value := range each.Headers {
						next.Add(name, value)
					}
				}
			}
		}
	}
}
