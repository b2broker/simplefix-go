package session

import (
	"errors"
	simplefixgo "github.com/b2broker/simplefix-go"
)

var (
	ErrNotEnoughMessages = errors.New("not enough messages in the storage")
	ErrInvalidBoundaries = errors.New("invalid boundaries")
	ErrInvalidSequence   = errors.New("unexpected sequence index")
)

// MessageStorage is an interface providing a basic method for storing messages awaiting to be sent.
type MessageStorage interface {
	Save(pk string, msg simplefixgo.SendingMessage, msgSeqNum int) error
	Messages(pk string, msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error)
}
