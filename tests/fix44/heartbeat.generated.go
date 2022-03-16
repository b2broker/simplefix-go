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

func CreateHeartbeat() *Heartbeat {
	msg := makeHeartbeat()

	return msg
}

func NewHeartbeat() *Heartbeat {
	m := makeHeartbeat()
	return &Heartbeat{
		fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(m.Body()...).
			SetHeader(m.Header().AsComponent()).
			SetTrailer(m.Trailer().AsComponent()),
	}
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

func (heartbeat *Heartbeat) SetFieldTestReqID(testReqID string) messages.HeartbeatBuilder {
	return heartbeat.SetTestReqID(testReqID)
}
