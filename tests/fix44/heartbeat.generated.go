package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeHeartbeat = "0"

type Heartbeat struct {
	*fix.Message
}

func makeHeartbeat() *Heartbeat {
	msg := &Heartbeat{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(
				fix.NewKeyValue(FieldTestReqID, &fix.String{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func NewHeartbeat() *Heartbeat {
	msg := makeHeartbeat()

	return msg
}

func ParseHeartbeat(data []byte) (*Heartbeat, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, FieldBeginString, beginString).
		SetBody(makeHeartbeat().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &Heartbeat{
		Message: msg,
	}, nil
}

func (heartbeat *Heartbeat) Header() *Header {
	header := heartbeat.Message.Header()

	return &Header{header}
}

func (heartbeat *Heartbeat) HeaderBuilder() messages.HeaderBuilder {
	return heartbeat.Header()
}

func (heartbeat *Heartbeat) Trailer() *Trailer {
	trailer := heartbeat.Message.Trailer()

	return &Trailer{trailer}
}

func (heartbeat *Heartbeat) TestReqID() string {
	kv := heartbeat.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (heartbeat *Heartbeat) SetTestReqID(testReqID string) *Heartbeat {
	kv := heartbeat.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(testReqID)
	return heartbeat
}

func (Heartbeat) New() messages.HeartbeatBuilder {
	return makeHeartbeat()
}

func (Heartbeat) Parse(data []byte) (messages.HeartbeatBuilder, error) {
	return ParseHeartbeat(data)
}

func (heartbeat *Heartbeat) SetFieldTestReqID(testReqID string) messages.HeartbeatBuilder {
	return heartbeat.SetTestReqID(testReqID)
}
