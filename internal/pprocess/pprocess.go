// Package pprocess stands for Periodic Process
package pprocess

import "context"

// PProcess defined the periodic process interface
type PProcess interface {
	Start(context.Context) error
	Stop() error
}
