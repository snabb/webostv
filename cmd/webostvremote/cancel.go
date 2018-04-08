package main

import (
	"sync"
)

type CancelPrevious struct {
	ch chan struct{}
	sync.Mutex
}

func (ce *CancelPrevious) NewCancel() (ch chan struct{}) {
	ch = make(chan struct{})

	ce.Lock()
	if ce.ch != nil {
		close(ce.ch)
	}
	ce.ch = ch
	ce.Unlock()

	return ch
}

func (ce *CancelPrevious) Cancel() {
	ce.Lock()
	if ce.ch != nil {
		close(ce.ch)
	}
	ce.ch = nil
	ce.Unlock()
}
