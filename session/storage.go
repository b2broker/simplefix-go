package session

import (
	"errors"
	simplefixgo "github.com/b2broker/simplefix-go"
)

var (
	ErrNotEnoughMessages = errors.New("not enough messages in storage")
	ErrInvalidBoundaries = errors.New("invalid boundaries")
	ErrInvalidSequence   = errors.New("unexpected sequence index")
)

// MessageStorage is an interface with basic method for storage of sending messages
type MessageStorage interface {
	Save(msg simplefixgo.SendingMessage, msgSeqNum int) error
	Messages(msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error)
}
