package md2json

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var JsonDocsPool sync.Pool

type StringSlicePool struct {
	Pool   sync.Pool
	Inited atomic.Bool
}

func NewStringSlicePool(capacity int) (ssp *StringSlicePool) {
	ssp = &StringSlicePool{
		Pool: sync.Pool{
			New: func() interface{} {
				return make([]string, 0, capacity)
			},
		},
		Inited: atomic.Bool{},
	}

	ssp.Inited.Store(true)

	return
}

func (ssp *StringSlicePool) Get() (slice []string, err error) {
	if ssp.Inited.Load() {
		slice = ssp.Pool.Get().([]string)
		return
	}
	err = fmt.Errorf("StringSlicePool is not initialized")
	return
}

func (ssp *StringSlicePool) Put(strSlice []string) {
	if ssp.Inited.Load() {
		ssp.Pool.Put(strSlice[:0])
	}
}
