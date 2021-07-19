package memory

import (
	"bytes"
	"errors"
	"github.com/b2broker/simplefix-go/session"
	"testing"
)

func TestStorage_Save(t *testing.T) {
	data := [][]byte{
		{}, {}, {}, {},
	}

	st := NewStorage(10, 5)
	for num, msg := range data {
		err := st.Save(msg, num+1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	err := st.Save([]byte{}, 10)
	if !errors.Is(err, session.ErrInvalidSequence) {
		t.Fatalf("expect error '%s', got '%s'", session.ErrInvalidSequence, err)
	}
}

func TestStorage_Messages(t *testing.T) {
	data := [][]byte{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		[]byte("test_me_11"),
		[]byte("test_me_12"),
		{}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	st := NewStorage(10, 5)
	for num, msg := range data {
		err := st.Save(msg, num+1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	messages, err := st.Messages(11, 12)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(messages) != 2 {
		t.Fatalf("unexpected size of messages: got %d, expect: %d", len(messages), 2)
	}

	if !bytes.Equal(messages[0], data[10]) {
		t.Fatalf("unexpected message: got %s, expect: %s", messages[0], data[10])
	}
	if !bytes.Equal(messages[1], data[11]) {
		t.Fatalf("unexpected message: got %s, expect: %s", messages[0], data[11])
	}

	_, err = st.Messages(100, 10)
	if !errors.Is(err, session.ErrInvalidBoundaries) {
		t.Fatalf("expect error '%s', got '%s'", session.ErrInvalidBoundaries, err)
	}
	_, err = st.Messages(5, 9)
	if !errors.Is(err, session.ErrNotEnoughMessages) {
		t.Fatalf("expect error '%s', got '%s'", session.ErrInvalidBoundaries, err)
	}
}
