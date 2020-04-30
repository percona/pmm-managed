package cleaner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/internal/pprocess"
)

// Cleaner holds unexported field of the background cleaner process.
type Cleaner struct {
	db       *reform.DB
	ticker   chan time.Time
	stopChan chan struct{}
	lock     *sync.Mutex
	running  bool
}

// New returns a new instance of the cleaner object that implements the PProcess interface.
func New(db *reform.DB, ticker chan time.Time) pprocess.PProcess {
	return &Cleaner{
		db:       db,
		ticker:   ticker,
		stopChan: make(chan struct{}),
		lock:     &sync.Mutex{},
	}
}

// Start starts the background process cleaner.
func (p *Cleaner) Start(ctx context.Context) error {
	if p.isRunning() {
		return fmt.Errorf("cleaner process already started")
	}
	go p.run(ctx)
	return nil
}

// Stop the background process.
func (p *Cleaner) Stop() error {
	if !p.isRunning() {
		return fmt.Errorf("cleaner process is not running")
	}

	p.stopChan <- struct{}{}
	return nil
}

func (p *Cleaner) run(ctx context.Context) {
	defer p.setRunning(false)

	// Semaphore to prevent having 2 queries running at the same time.
	// It should never happen but if for some reason the clean query takes too much time
	// and we receive a ticker tick while the query is still running, we should prevent that
	sem := make(chan struct{}, 1)
	sem <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-p.ticker:
			select {
			case <-sem:
				// TODO run the query to delete old data
				sem <- struct{}{} // restore the semaphore
			default:
				// the previous query is still running. skip this tick
			}
		}
	}
}

func (p *Cleaner) isRunning() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.running
}

func (p *Cleaner) setRunning(status bool) {
	p.lock.Lock()
	p.running = status
	p.lock.Unlock()
}
