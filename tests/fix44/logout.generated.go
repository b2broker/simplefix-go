package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeLogout = "5"

type Logout struct {
	*fix.Message
}

func makeLogout() *Logout {
	msg := &Logout{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeLogout).
			SetBody(
				fix.NewKeyValue(FieldText, &fix.String{}),
				fix.NewKeyValue(FieldEncodedTextLen, &fix.Int{}),
				fix.NewKeyValue(FieldEncodedText, &fix.String{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func CreateLogout() *Logout {
	msg := makeLogout()

	return msg
}

func NewLogout() *Logout {
	m := makeLogout()
	return &Logout{
		fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(m.Body()...).
			SetHeader(m.Header().AsComponent()).
			SetTrailer(m.Trailer().AsComponent()),
	}
}

func (logout *Logout) Header() *Header {
	header := logout.Message.Header()

	return &Header{header}
}

func (logout *Logout) HeaderBuilder() messages.HeaderBuilder {
	return logout.Header()
}

func (logout *Logout) Trailer() *Trailer {
	trailer := logout.Message.Trailer()

	return &Trailer{trailer}
}

func (logout *Logout) Text() string {
	kv := logout.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (logout *Logout) SetText(text string) *Logout {
	kv := logout.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(text)
	return logout
}

func (logout *Logout) EncodedTextLen() int {
	kv := logout.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (logout *Logout) SetEncodedTextLen(encodedTextLen int) *Logout {
	kv := logout.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(encodedTextLen)
	return logout
}

func (logout *Logout) EncodedText() string {
	kv := logout.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (logout *Logout) SetEncodedText(encodedText string) *Logout {
	kv := logout.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(encodedText)
	return logout
}

// New is a plane message constructor
func (Logout) New() messages.LogoutBuilder {
	return makeLogout()
}

// Build provides an opportunity to customize message during building outgoing message
func (Logout) Build() messages.LogoutBuilder {
	return makeLogout()
}
