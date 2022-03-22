package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeSequenceReset = "4"

type SequenceReset struct {
	*fix.Message
}

func makeSequenceReset() *SequenceReset {
	msg := &SequenceReset{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeSequenceReset).
			SetBody(
				fix.NewKeyValue(FieldGapFillFlag, &fix.String{}),
				fix.NewKeyValue(FieldNewSeqNo, &fix.Int{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func CreateSequenceReset(newSeqNo int) *SequenceReset {
	msg := makeSequenceReset().
		SetNewSeqNo(newSeqNo)

	return msg
}

func NewSequenceReset() *SequenceReset {
	m := makeSequenceReset()
	return &SequenceReset{
		fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeHeartbeat).
			SetBody(m.Body()...).
			SetHeader(m.Header().AsComponent()).
			SetTrailer(m.Trailer().AsComponent()),
	}
}

func (sequenceReset *SequenceReset) Header() *Header {
	header := sequenceReset.Message.Header()

	return &Header{header}
}

func (sequenceReset *SequenceReset) HeaderBuilder() messages.HeaderBuilder {
	return sequenceReset.Header()
}

func (sequenceReset *SequenceReset) Trailer() *Trailer {
	trailer := sequenceReset.Message.Trailer()

	return &Trailer{trailer}
}

func (sequenceReset *SequenceReset) GapFillFlag() string {
	kv := sequenceReset.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (sequenceReset *SequenceReset) SetGapFillFlag(gapFillFlag string) *SequenceReset {
	kv := sequenceReset.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(gapFillFlag)
	return sequenceReset
}

func (sequenceReset *SequenceReset) NewSeqNo() int {
	kv := sequenceReset.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (sequenceReset *SequenceReset) SetNewSeqNo(newSeqNo int) *SequenceReset {
	kv := sequenceReset.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(newSeqNo)
	return sequenceReset
}

// New is a plane message constructor
func (SequenceReset) New() messages.SequenceResetBuilder {
	return makeSequenceReset()
}

// Build provides an opportunity to customize message during building outgoing message
func (SequenceReset) Build() messages.SequenceResetBuilder {
	return makeSequenceReset()
}

func (sequenceReset *SequenceReset) SetFieldNewSeqNo(newSeqNo int) messages.SequenceResetBuilder {
	return sequenceReset.SetNewSeqNo(newSeqNo)
}
