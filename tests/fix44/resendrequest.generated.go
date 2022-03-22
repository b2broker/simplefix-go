package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeResendRequest = "2"

type ResendRequest struct {
	*fix.Message
}

func makeResendRequest() *ResendRequest {
	msg := &ResendRequest{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeResendRequest).
			SetBody(
				fix.NewKeyValue(FieldBeginSeqNo, &fix.Int{}),
				fix.NewKeyValue(FieldEndSeqNo, &fix.Int{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func CreateResendRequest(beginSeqNo int, endSeqNo int) *ResendRequest {
	msg := makeResendRequest().
		SetBeginSeqNo(beginSeqNo).
		SetEndSeqNo(endSeqNo)

	return msg
}

func NewResendRequest() *ResendRequest {
	m := makeResendRequest()
	return &ResendRequest{
		fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(m.Body()...).
			SetHeader(m.Header().AsComponent()).
			SetTrailer(m.Trailer().AsComponent()),
	}
}

func (resendRequest *ResendRequest) Header() *Header {
	header := resendRequest.Message.Header()

	return &Header{header}
}

func (resendRequest *ResendRequest) HeaderBuilder() messages.HeaderBuilder {
	return resendRequest.Header()
}

func (resendRequest *ResendRequest) Trailer() *Trailer {
	trailer := resendRequest.Message.Trailer()

	return &Trailer{trailer}
}

func (resendRequest *ResendRequest) BeginSeqNo() int {
	kv := resendRequest.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (resendRequest *ResendRequest) SetBeginSeqNo(beginSeqNo int) *ResendRequest {
	kv := resendRequest.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(beginSeqNo)
	return resendRequest
}

func (resendRequest *ResendRequest) EndSeqNo() int {
	kv := resendRequest.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (resendRequest *ResendRequest) SetEndSeqNo(endSeqNo int) *ResendRequest {
	kv := resendRequest.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(endSeqNo)
	return resendRequest
}

// New is a plane message constructor
func (ResendRequest) New() messages.ResendRequestBuilder {
	return makeResendRequest()
}

// Build provides an opportunity to customize message during building outgoing message
func (ResendRequest) Build() messages.ResendRequestBuilder {
	return makeResendRequest()
}

func (resendRequest *ResendRequest) SetFieldBeginSeqNo(beginSeqNo int) messages.ResendRequestBuilder {
	return resendRequest.SetBeginSeqNo(beginSeqNo)
}

func (resendRequest *ResendRequest) SetFieldEndSeqNo(endSeqNo int) messages.ResendRequestBuilder {
	return resendRequest.SetEndSeqNo(endSeqNo)
}
