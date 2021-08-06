package utils

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrZeroTimeout       = errors.New("zero timeout")
	ErrTooSmallFrequency = errors.New("too small frequency")
)

const frequency time.Duration = 10
const minFrequency = time.Microsecond

type Timer struct {
	mu         sync.RWMutex
	lastUpdate time.Time

	timeout         time.Duration
	checkingTimeout time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimer(timeout time.Duration) (*Timer, error) {
	if timeout == 0 {
		return nil, ErrZeroTimeout
	}

	checkingTimeout := timeout / frequency

	if checkingTimeout < minFrequency {
		return nil, ErrTooSmallFrequency
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Timer{
		checkingTimeout: checkingTimeout,
		timeout:         timeout,
		ctx:             ctx,
		cancel:          cancel,
	}, nil
}

func (t *Timer) Close() {
	t.cancel()
}

func (t *Timer) Refresh() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastUpdate = time.Now()
}

// TakeTimeout will be end when time will is over or closed
func (t *Timer) TakeTimeout() {
	t.Refresh()
	ticker := time.NewTicker(t.checkingTimeout)
	for {
		select {
		case <-ticker.C:
			t.mu.RLock()
			rest := time.Until(t.lastUpdate.Add(t.timeout))
			t.mu.RUnlock()

			if rest <= 0 {
				return
			}

		case <-t.ctx.Done():
			return
		}
	}
}
