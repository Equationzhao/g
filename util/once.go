package util

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	m    sync.RWMutex
	done uint32
}

func (o *Once) Do(fn func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.doSlow(fn)
}

func (o *Once) doSlow(fn func() error) error {
	o.m.RLock()
	var err error
	if o.done == 0 {
		o.m.RUnlock()
		o.m.Lock()

		defer o.m.Unlock()
		if o.done == 0 {
			err = fn()
			if err == nil {
				atomic.StoreUint32(&o.done, 1)
			}
		}

	} else {
		o.m.RUnlock()
	}
	return err
}
