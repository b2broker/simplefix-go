package messages

type MockMessage struct {
	Type string
	Data []byte
	Err  error
}

func NewMockMessage(tp string, data []byte, err error) *MockMessage {
	return &MockMessage{Type: tp, Data: data, Err: err}
}

func (m MockMessage) HeaderBuilder() HeaderBuilder {
	return nil
}

func (m MockMessage) MsgType() string {
	return m.Type
}

func (m MockMessage) ToBytes() ([]byte, error) {
	return m.Data, m.Err
}
