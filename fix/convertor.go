package fix

import (
	"sync"

	"github.com/b2broker/simplefix-go/fix/buffer"
)

type MessageByteConverter struct {
	pool sync.Pool
}

func NewMessageByteConverter(bufferSize int) *MessageByteConverter {
	b := &MessageByteConverter{
		pool: sync.Pool{
			New: func() interface{} {
				return buffer.NewMessageByteBuffers(bufferSize)
			},
		},
	}
	return b
}

type ConvertableMessage interface {
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
}

func (m *MessageByteConverter) ConvertToBytes(message ConvertableMessage) ([]byte, error) {
	buffers := m.pool.Get().(*buffer.MessageByteBuffers)
	buffers.Reset()
	defer m.pool.Put(buffers)

	return message.ToBytesBuffered(buffers)
}
