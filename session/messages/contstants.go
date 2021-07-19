package messages

type Tags struct {
	MsgType         int
	MsgSeqNum       int
	HeartBtInt      int
	EncryptedMethod int
}

type SessionErrorCodes struct {
	RequiredTagMissing int
	IncorrectValue     int
	Other              int
}

type Message interface {
	HeaderBuilder() HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
}
