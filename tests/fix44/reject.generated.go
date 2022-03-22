package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeReject = "3"

type Reject struct {
	*fix.Message
}

func makeReject() *Reject {
	msg := &Reject{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeReject).
			SetBody(
				fix.NewKeyValue(FieldRefSeqNum, &fix.Int{}),
				fix.NewKeyValue(FieldRefTagID, &fix.Int{}),
				fix.NewKeyValue(FieldRefMsgType, &fix.String{}),
				fix.NewKeyValue(FieldSessionRejectReason, &fix.String{}),
				fix.NewKeyValue(FieldText, &fix.String{}),
				fix.NewKeyValue(FieldEncodedTextLen, &fix.Int{}),
				fix.NewKeyValue(FieldEncodedText, &fix.String{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func CreateReject(refSeqNum int) *Reject {
	msg := makeReject().
		SetRefSeqNum(refSeqNum)

	return msg
}

func NewReject() *Reject {
	m := makeReject()
	return &Reject{
		fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(m.Body()...).
			SetHeader(m.Header().AsComponent()).
			SetTrailer(m.Trailer().AsComponent()),
	}
}

func (reject *Reject) Header() *Header {
	header := reject.Message.Header()

	return &Header{header}
}

func (reject *Reject) HeaderBuilder() messages.HeaderBuilder {
	return reject.Header()
}

func (reject *Reject) Trailer() *Trailer {
	trailer := reject.Message.Trailer()

	return &Trailer{trailer}
}

func (reject *Reject) RefSeqNum() int {
	kv := reject.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (reject *Reject) SetRefSeqNum(refSeqNum int) *Reject {
	kv := reject.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(refSeqNum)
	return reject
}

func (reject *Reject) RefTagID() int {
	kv := reject.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (reject *Reject) SetRefTagID(refTagID int) *Reject {
	kv := reject.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(refTagID)
	return reject
}

func (reject *Reject) RefMsgType() string {
	kv := reject.Get(2)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (reject *Reject) SetRefMsgType(refMsgType string) *Reject {
	kv := reject.Get(2).(*fix.KeyValue)
	_ = kv.Load().Set(refMsgType)
	return reject
}

func (reject *Reject) SessionRejectReason() string {
	kv := reject.Get(3)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (reject *Reject) SetSessionRejectReason(sessionRejectReason string) *Reject {
	kv := reject.Get(3).(*fix.KeyValue)
	_ = kv.Load().Set(sessionRejectReason)
	return reject
}

func (reject *Reject) Text() string {
	kv := reject.Get(4)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (reject *Reject) SetText(text string) *Reject {
	kv := reject.Get(4).(*fix.KeyValue)
	_ = kv.Load().Set(text)
	return reject
}

func (reject *Reject) EncodedTextLen() int {
	kv := reject.Get(5)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (reject *Reject) SetEncodedTextLen(encodedTextLen int) *Reject {
	kv := reject.Get(5).(*fix.KeyValue)
	_ = kv.Load().Set(encodedTextLen)
	return reject
}

func (reject *Reject) EncodedText() string {
	kv := reject.Get(6)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (reject *Reject) SetEncodedText(encodedText string) *Reject {
	kv := reject.Get(6).(*fix.KeyValue)
	_ = kv.Load().Set(encodedText)
	return reject
}

// New is a plane message constructor
func (Reject) New() messages.RejectBuilder {
	return makeReject()
}

// Build provides an opportunity to customize message during building outgoing message
func (Reject) Build() messages.RejectBuilder {
	return makeReject()
}

func (reject *Reject) SetFieldSessionRejectReason(sessionRejectReason string) messages.RejectBuilder {
	return reject.SetSessionRejectReason(sessionRejectReason)
}

func (reject *Reject) SetFieldRefSeqNum(refSeqNum int) messages.RejectBuilder {
	return reject.SetRefSeqNum(refSeqNum)
}

func (reject *Reject) SetFieldRefTagID(refTagID int) messages.RejectBuilder {
	return reject.SetRefTagID(refTagID)
}
