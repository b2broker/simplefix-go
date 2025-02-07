package messages

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/fix/buffer"
)

type Builder interface {
	Items() fix.Items
	CalcBodyLength() int
	BodyLength() int
	BytesWithoutChecksum() []byte
	CheckSum() string
	BeginString() *fix.KeyValue
	MsgType() string
	ToBytes() ([]byte, error)
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
	BeginStringTag() string
	BodyLengthTag() string
	CheckSumTag() string
}

type PipelineBuilder interface {
	HeaderBuilder() HeaderBuilder
	Builder
}
