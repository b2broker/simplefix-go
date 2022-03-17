package messages

import "github.com/b2broker/simplefix-go/fix"

type Builder interface {
	Items() fix.Items
	CalcBodyLength() int
	BodyLength() int
	BytesWithoutChecksum() []byte
	CheckSum() string
	BeginString() *fix.KeyValue
	MsgType() string
	ToBytes() ([]byte, error)
}

type PipelineBuilder interface {
	HeaderBuilder() HeaderBuilder
	Builder
}
