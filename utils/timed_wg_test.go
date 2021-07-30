package utils

import (
	"errors"
	"testing"
	"time"
)

func TestTimedWaitGroup_WaitWithTimeoutNegative(t *testing.T) {
	wg := TimedWaitGroup{}

	wg.Add(1)

	go func() {
		time.Sleep(time.Millisecond * 10)
		wg.Done()
	}()

	err := wg.WaitWithTimeout(time.Millisecond)
	if !errors.Is(err, ErrWGExpired) {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestTimedWaitGroup_WaitWithTimeout(t *testing.T) {
	wg := TimedWaitGroup{}

	wg.Add(1)

	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
	}()

	err := wg.WaitWithTimeout(time.Millisecond * 2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
