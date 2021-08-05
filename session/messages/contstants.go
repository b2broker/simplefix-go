package messages

// Tags required tags for session pipelines
type Tags struct {
	MsgType         int
	MsgSeqNum       int
	HeartBtInt      int
	EncryptedMethod int
}

// SessionErrorCodes session error codes
type SessionErrorCodes struct {
	RequiredTagMissing int
	IncorrectValue     int
	Other              int
}

// Message interface for sending message
type Message interface {
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
