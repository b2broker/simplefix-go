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

func NewResendRequest(beginSeqNo int, endSeqNo int) *ResendRequest {
	msg := makeResendRequest().
		SetBeginSeqNo(beginSeqNo).
		SetEndSeqNo(endSeqNo)

	return msg
}

func ParseResendRequest(data []byte) (*ResendRequest, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, FieldBeginString, beginString).
		SetBody(makeResendRequest().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &ResendRequest{
		Message: msg,
	}, nil
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

func (ResendRequest) New() messages.ResendRequestBuilder {
	return makeResendRequest()
}

func (ResendRequest) Parse(data []byte) (messages.ResendRequestBuilder, error) {
	return ParseResendRequest(data)
}

func (resendRequest *ResendRequest) SetFieldBeginSeqNo(beginSeqNo int) messages.ResendRequestBuilder {
	return resendRequest.SetBeginSeqNo(beginSeqNo)
}

func (resendRequest *ResendRequest) SetFieldEndSeqNo(endSeqNo int) messages.ResendRequestBuilder {
	return resendRequest.SetEndSeqNo(endSeqNo)
}
