package pmmapitests

import (
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var Context context.Context
var BaseURL *url.URL

func init() {
	debugF := flag.Bool("pmm.debug", false, "Enable debug output.")
	serverURLF := flag.String("pmm.server-url", "https://127.0.0.1:8443/", "PMM Server URL.")
	flag.Parse()

	if *debugF {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled.")
	}

	var cancel func()
	Context, cancel = context.WithCancel(context.Background())

	// handle termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-signals
		signal.Stop(signals)
		logrus.Warnf("Got %s, shutting down...", unix.SignalName(s.(syscall.Signal)))
		cancel()
	}()

	var err error
	BaseURL, err = url.Parse(*serverURLF)
	if err != nil {
		logrus.Fatalf("Failed to parse PMM Server URL: %s.", err)
	}
	if BaseURL.Host == "" || BaseURL.Scheme == "" {
		logrus.Fatalf("Invalid PMM Server URL: %s", BaseURL.String())
	}
	if BaseURL.Path == "" {
		BaseURL.Path = "/"
	}
	logrus.Debugf("PMM Server URL: %#v.", BaseURL)

	// use JSON APIs over HTTP/1.1 (setting TLSNextProto to non-nil map disables automated HTTP/2)
	transport := httptransport.New(BaseURL.Host, BaseURL.Path, []string{BaseURL.Scheme})
	transport.Transport.(*http.Transport).TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}

	l := logrus.WithField("component", "client")
	transport.SetLogger(l)
	transport.Debug = *debugF

	client.Default = client.New(transport, nil)
}
