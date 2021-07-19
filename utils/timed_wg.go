package utils

import (
	"errors"
	"sync"
	"time"
)

var ErrWGExpired = errors.New("wait group timeout expired")

type TimedWaitGroup struct {
	sync.WaitGroup
}

func (wg *TimedWaitGroup) WaitWithTimeout(timeout time.Duration) error {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ch:
		return nil

	case <-time.After(timeout):
		return ErrWGExpired
	}
}
