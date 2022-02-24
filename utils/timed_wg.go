package utils

import (
	"errors"
	"sync"
	"time"
)

// ErrWGExpired is returned if the timeout expires before the WaitGroup has received
// a "Done" confirmation from all of the streams.
var ErrWGExpired = errors.New("wait group timeout expired")

// TimedWaitGroup is a combination of WaitGroup and timeout
type TimedWaitGroup struct {
	sync.WaitGroup
}

// WaitWithTimeout awaits any of the two cases:
// - timeout expires (in which case an error is returned)
// - a WaitGroup receives a "Done" confirmation (in which case nil is returned)
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
