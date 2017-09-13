package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestRunTelemetry(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	logrus.SetOutput(w)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	ch := make(chan bool)
	go func() {
		runTelemetryService(ctx)
		ch <- true
	}()

	select {
	case <-ch:
		break
	case <-time.After(2100 * time.Millisecond):
		t.Fatalf("Service didn't stop upon context cancellation")
	}

	w.Flush()
	assert.Contains(t, b.String(), `msg="Telemetry is enabled. Send data interval = 24h0m0s" component=TELEMETRY`)

	logrus.SetOutput(os.Stdout)

}
