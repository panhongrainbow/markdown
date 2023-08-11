package md2json

import (
	"sync"
	"sync/atomic"
)

// StringSlicePool represents a pool of string slices.
type StringSlicePool struct {
	Pool   sync.Pool
	Inited atomic.Bool
}

// NewStringSlicePool creates a new StringSlicePool with the specified capacity.
func NewStringSlicePool(capacity int) (ssp *StringSlicePool) {
	ssp = &StringSlicePool{
		Pool: sync.Pool{
			// Create a new sync.Pool to hold string slices.
			// The New function initializes each new string slice with a capacity.
			New: func() interface{} {
				return make([]string, 0, capacity)
			},
		},
		// Initialize the inited flag as false.
		Inited: atomic.Bool{},
	}

	// Mark the pool as initialized.
	ssp.Inited.Store(true)

	return
}

// Get retrieves a string slice from the pool.
func (ssp *StringSlicePool) Get() (slice []string) {
	if ssp.Inited.Load() {
		slice = ssp.Pool.Get().([]string)
		return
	}

	return
}

// Put returns a string slice to the pool.
func (ssp *StringSlicePool) Put(strSlice []string) {
	if ssp.Inited.Load() {
		ssp.Pool.Put(strSlice[:0])
	}
}
