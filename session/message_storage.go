package session

import (
	simplefixgo "github.com/b2broker/simplefix-go"
)

// MessageStorage is an interface providing a basic method for storing messages awaiting to be sent.
type MessageStorage interface {
	Save(pk string, msg simplefixgo.SendingMessage, msgSeqNum int) error
	Messages(pk string, msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error)
}
