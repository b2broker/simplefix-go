package memory

import (
	"errors"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	"testing"
)

func TestStorage_Save(t *testing.T) {
	data := []simplefixgo.SendingMessage{
		messages.NewMockMessage("12", []byte{}, nil),
		messages.NewMockMessage("12", []byte{}, nil),
		messages.NewMockMessage("12", []byte{}, nil),
		messages.NewMockMessage("12", []byte{}, nil),
	}

	st := NewStorage(10, 5)
	for num, msg := range data {
		err := st.Save(msg, num+1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	err := st.Save(nil, 10)
	if !errors.Is(err, session.ErrInvalidSequence) {
		t.Fatalf("expect error '%s', got '%s'", session.ErrInvalidSequence, err)
	}
}

func TestStorage_Messages(t *testing.T) {
	data := []simplefixgo.SendingMessage{
		messages.NewMockMessage("01", []byte{}, nil),
		messages.NewMockMessage("02", []byte{}, nil),
		messages.NewMockMessage("03", []byte{}, nil),
		messages.NewMockMessage("04", []byte{}, nil),
		messages.NewMockMessage("05", []byte{}, nil),
		messages.NewMockMessage("06", []byte{}, nil),
		messages.NewMockMessage("07", []byte{}, nil),
		messages.NewMockMessage("08", []byte{}, nil),
		messages.NewMockMessage("09", []byte{}, nil),
		messages.NewMockMessage("10", []byte{}, nil),
		messages.NewMockMessage("11", []byte{}, nil),
		messages.NewMockMessage("12", []byte{}, nil),
		messages.NewMockMessage("13", []byte{}, nil),
		messages.NewMockMessage("14", []byte{}, nil),
		messages.NewMockMessage("15", []byte{}, nil),
		messages.NewMockMessage("16", []byte{}, nil),
		messages.NewMockMessage("17", []byte{}, nil),
		messages.NewMockMessage("18", []byte{}, nil),
		messages.NewMockMessage("19", []byte{}, nil),
	}

	st := NewStorage(10, 5)
	for num, msg := range data {
		err := st.Save(msg, num+1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	storageMessages, err := st.Messages(11, 12)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(storageMessages) != 2 {
		t.Fatalf("unexpected size of storageMessages: got %d, expect: %d", len(storageMessages), 2)
	}

	if storageMessages[0].MsgType() == "" || storageMessages[0].MsgType() != data[10].MsgType() {
		t.Fatalf("unexpected message: got %s, expect: %s", storageMessages[0], data[10])
	}

	if storageMessages[1].MsgType() == "" || storageMessages[1].MsgType() != data[11].MsgType() {
		t.Fatalf("unexpected message: got %s, expect: %s", storageMessages[0], data[11])
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
