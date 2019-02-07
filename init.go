package pmmapitests

import (
	"crypto/tls"
	"flag"
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/sirupsen/logrus"
)

func init() {
	serverURLF := flag.String("pmm.server-url", "https://127.0.0.1:8443/", "PMM Server URL.")
	debugF := flag.Bool("pmm.debug", false, "Enable debug output.")
	flag.Parse()

	if *debugF {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled.")
	}

	u, err := url.Parse(*serverURLF)
	if err != nil {
		logrus.Fatalf("Failed to parse PMM Server URL: %s.", err)
	}
	if u.Host == "" || u.Scheme == "" {
		logrus.Fatalf("Invalid PMM Server URL: %s", u.String())
	}
	if u.Path == "" {
		u.Path = "/"
	}
	logrus.Debugf("PMM Server URL: %#v.", u)

	// use JSON APIs over HTTP/1.1 (setting TLSNextProto to non-nil map disables automated HTTP/2)
	transport := httptransport.New(u.Host, u.Path, []string{u.Scheme})
	transport.Transport.(*http.Transport).TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}

	l := logrus.WithField("component", "client")
	transport.SetLogger(l)
	transport.Debug = *debugF

	client.Default = client.New(transport, nil)
}
