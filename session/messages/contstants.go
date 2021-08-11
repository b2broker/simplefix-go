package messages

// Tags contains required tags for session pipelines
type Tags struct {
	MsgType         int
	MsgSeqNum       int
	HeartBtInt      int
	EncryptedMethod int
}

// SessionErrorCodes contains session error codes
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

// Message is an interface for sending message
type Message interface {
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
