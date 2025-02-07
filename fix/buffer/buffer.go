package buffer

import "bytes"

type MessageByteBuffers struct {
	msgBuffer    *bytes.Buffer
	bodyBuffer   *bytes.Buffer
	typeBuffer   *bytes.Buffer
	headerBuffer *bytes.Buffer
}

func NewMessageByteBuffers(size int) *MessageByteBuffers {
	return &MessageByteBuffers{
		typeBuffer:   bytes.NewBuffer(make([]byte, 0, 5)), // 35=AA
		msgBuffer:    bytes.NewBuffer(make([]byte, 0, size)),
		bodyBuffer:   bytes.NewBuffer(make([]byte, 0, size)),
		headerBuffer: bytes.NewBuffer(make([]byte, 0, size)),
	}
}

func (m *MessageByteBuffers) Reset() {
	m.msgBuffer.Reset()
	m.bodyBuffer.Reset()
	m.typeBuffer.Reset()
	m.headerBuffer.Reset()
}

func (m *MessageByteBuffers) GetMessageBuffer() *bytes.Buffer {
	return m.msgBuffer
}
func (m *MessageByteBuffers) GetBodyBuffer() *bytes.Buffer {
	return m.bodyBuffer
}
func (m *MessageByteBuffers) GetHeaderBuffer() *bytes.Buffer {
	return m.headerBuffer
}
func (m *MessageByteBuffers) GetTypeBuffer() *bytes.Buffer {
	return m.typeBuffer
}
