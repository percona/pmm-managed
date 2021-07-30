// This file contains implementation for concurrent safe RNG.
package pmmapitests

import (
	"math/rand"
	"sync"
)

// ConcurrentRand wraps rand.Rand with mutex.
type ConcurrentRand struct {
	m    sync.Mutex
	rand *rand.Rand
}

// NewConcurrentRand constructs new ConcurrentRand with provided seed.
func NewConcurrentRand(seed int64) *ConcurrentRand {
	r := &ConcurrentRand{
		rand: rand.New(rand.NewSource(seed)),
	}
	return r
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
func (r *ConcurrentRand) Seed(seed int64) {
	r.m.Lock()
	defer r.m.Unlock()
	r.rand.Seed(seed)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *ConcurrentRand) Int63() int64 {
	r.m.Lock()
	defer r.m.Unlock()
	return r.rand.Int63()
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (r *ConcurrentRand) Uint64() uint64 {
	r.m.Lock()
	defer r.m.Unlock()
	return r.rand.Uint64()
}
