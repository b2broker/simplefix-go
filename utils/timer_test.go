package utils

import (
	"errors"
	"testing"
	"time"
)

func TestTimer_TimeoutZero(t *testing.T) {
	_, err := NewTimer(0)
	if !errors.Is(err, ErrZeroTimeout) {
		t.Fatalf("expected error: %s, returned: %s", ErrZeroTimeout, err)
	}
}

func TestTimer_Timeout(t *testing.T) {
	delay := time.Millisecond * 10
	tm, err := NewTimer(delay)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	now := time.Now()
	for {
		tm.TakeTimeout()
		if time.Until(now.Add(delay)) > 0 {
			t.Fatalf("delay was not applied")
		} else {
			return
		}
	}
}
