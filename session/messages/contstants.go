package messages

import (
	"github.com/b2broker/simplefix-go/fix/buffer"
)

// Tags is a structure specifying the required tags for session pipelines.
type Tags struct {
	MsgType         int
	MsgSeqNum       int
	HeartBtInt      int
	EncryptedMethod int
}

// SessionErrorCodes is a structure specifying the session error codes.
type SessionErrorCodes struct {
	InvalidTagNumber            int
	RequiredTagMissing          int
	TagNotDefinedForMessageType int
	UndefinedTag                int
	TagSpecialWithoutValue      int
	IncorrectValue              int
	IncorrectDataFormatValue    int
	DecryptionProblem           int
	SignatureProblem            int
	CompIDProblem               int
	Other                       int
}

// Message is an interface providing the functionality required for sending messages.
type Message interface {
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
}
