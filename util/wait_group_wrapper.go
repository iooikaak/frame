package util

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func (w *WaitGroupWrapper) WrapMuti(cb func(), goroutines int) {
	for i := 1; i <= goroutines; i++ {
		w.Add(i)
		go func() {
			cb()
			w.Done()
		}()
	}
}
