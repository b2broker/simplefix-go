package utils

import (
	"errors"
	"sync"
	"time"
)

// ErrWGExpired returns when timeout expired before wait group will be done
var ErrWGExpired = errors.New("wait group timeout expired")

type TimedWaitGroup struct {
	sync.WaitGroup
}

// WaitWithTimeout wait for one of two cases:
// timeout expired (returns error)
// wait group will be done (returns nil)
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
