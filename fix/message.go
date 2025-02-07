package fix

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/b2broker/simplefix-go/fix/buffer"
)

// Message is a structure providing functionality to FIX messages.
type Message struct {
	// The beginString, bodyLength and msgType fields are required for any FIX message.
	beginString *KeyValue
	bodyLength  *KeyValue
	msgType     *KeyValue // This must be the third tag appearing in the message header.

	// The message header must include all the required fields (bodyLength, beginString, msgType).
	header *Component

	// The message body.
	body Items

	// The message trailer which includes all the required fields except for checkSum.
	trailer *Component

	// The auto-generated checkSum value is a required field.
	checkSum *KeyValue

	prepared []byte
}

// NewMessage is called to create a new message.
func NewMessage(beginStringTag, bodyLengthTag, checkSumTag, msgTypeTag, beginString, msgType string) *Message {
	return &Message{
		beginString: NewKeyValue(beginStringTag, NewString(beginString)),
		bodyLength:  NewKeyValue(bodyLengthTag, &Int{}),
		msgType:     NewKeyValue(msgTypeTag, NewString(msgType)),
		checkSum:    NewKeyValue(checkSumTag, &String{}),
	}
}

// Items returns all fields as Items slice
func (msg *Message) Items() Items {
	items := Items{
		msg.beginString,
		msg.bodyLength,
		msg.msgType,
		msg.header,
	}

	items = append(items, msg.body...)
	items = append(items, msg.trailer, msg.checkSum)

	return items
}

// Body returns the body of a FIX message as an Items object.
func (msg *Message) Body() (kvs Items) {
	return msg.body
}

// Header returns the header of a FIX message as a Component object.
func (msg *Message) Header() *Component {
	return msg.header
}

// Trailer returns the trailer of a FIX message as a Component object.
func (msg *Message) Trailer() *Component {
	return msg.trailer
}

// BeginString returns the beginString field value.
func (msg *Message) BeginString() *KeyValue {
	return msg.beginString
}

func (msg *Message) BeginStringTag() string {
	return msg.beginString.Key
}

// BodyLength returns the BodyLength field value.
func (msg *Message) BodyLength() int {
	return msg.bodyLength.Value.Value().(int)
}

func (msg *Message) BodyLengthTag() string {
	return msg.bodyLength.Key
}

// MsgType returns the MsgType field value.
func (msg *Message) MsgType() string {
	return msg.msgType.Value.String()
}

// CheckSum returns the CheckSum field value.
func (msg *Message) CheckSum() string {
	return msg.checkSum.Value.String()
}

func (msg *Message) CheckSumTag() string {
	return msg.checkSum.Key
}

func (msg *Message) CalcBodyLength() int {
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()
	mt := msg.msgType.ToBytes()

	var length int

	if len(mt) > 0 {
		length += len(mt) + 1
	}

	if len(bb) > 0 {
		length += len(bb) + 1
	}

	if len(bh) > 0 {
		length += len(bh) + 1
	}

	return length
}

func (msg *Message) BytesWithoutChecksum() []byte {
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()

	bm := bytes.Join([][]byte{
		msg.beginString.ToBytes(),
		msg.bodyLength.ToBytes(),
		msg.msgType.ToBytes(),
	}, Delimiter)

	if len(bh) > 0 {
		bm = bytes.Join([][]byte{bm, bh}, Delimiter)
	}

	if len(bb) > 0 {
		bm = bytes.Join([][]byte{bm, bb}, Delimiter)
	}

	return bm
}

func (msg *Message) WriteBytesWithoutChecksum(buffers *buffer.MessageByteBuffers) {
	messageBuffer := buffers.GetMessageBuffer()
	bodyBuffer := buffers.GetBodyBuffer()
	headerBuffer := buffers.GetHeaderBuffer()

	bodyWritten := msg.body.WriteBytes(bodyBuffer)
	headerWritten := msg.header.WriteBytes(headerBuffer)
	totalLength := 0
	if headerWritten {
		totalLength += headerBuffer.Len() + 1
	}
	if bodyWritten {
		totalLength += bodyBuffer.Len() + 1
	}

	typeBuffer := buffers.GetTypeBuffer()
	if msg.msgType.WriteBytes(typeBuffer) {
		totalLength += typeBuffer.Len() + 1
	}

	msg.bodyLength.Value = NewInt(totalLength)

	msg.beginString.WriteBytes(messageBuffer)
	_ = messageBuffer.WriteByte(DelimiterChar)
	msg.bodyLength.WriteBytes(messageBuffer)
	_ = messageBuffer.WriteByte(DelimiterChar)
	msg.msgType.WriteBytes(messageBuffer)

	if headerWritten {
		_ = messageBuffer.WriteByte(DelimiterChar)
		_, _ = messageBuffer.Write(headerBuffer.Bytes())
	}

	if bodyWritten {
		_ = messageBuffer.WriteByte(DelimiterChar)
		_, _ = messageBuffer.Write(bodyBuffer.Bytes())
	}
}

// Prepare prepares message by calculating body length and check sum
func (msg *Message) Prepare() error {
	msg.bodyLength.Value = NewInt(msg.CalcBodyLength())

	byteMsg := msg.BytesWithoutChecksum()

	checkSum := CalcCheckSum(byteMsg)
	err := msg.checkSum.Value.Set(string(checkSum))
	if err != nil {
		return err
	}

	msg.prepared = bytes.Join([][]byte{
		byteMsg,
		makeTagValue(msg.checkSum.Key, checkSum),
	}, Delimiter)
	msg.prepared = append(msg.prepared, Delimiter...)

	return nil
}

func (msg *Message) prepareBuffered(buffers *buffer.MessageByteBuffers) error {
	msgBuff := buffers.GetMessageBuffer()

	msg.WriteBytesWithoutChecksum(buffers)
	checkSum := CalcCheckSumOptimizedFromBuffer(msgBuff)
	if err := msg.checkSum.Value.Set(string(checkSum)); err != nil {
		return err
	}

	msgBuff.WriteByte(DelimiterChar)
	_, _ = msgBuff.WriteString(msg.checkSum.Key)
	_ = msgBuff.WriteByte('=')
	_, _ = msgBuff.Write(checkSum)
	msgBuff.WriteByte(DelimiterChar)

	msg.prepared = make([]byte, msgBuff.Len())
	copy(msg.prepared, msgBuff.Bytes())

	return nil
}

// Prepared returns a byte representation of a specified message.
func (msg *Message) ToBytes() ([]byte, error) {

	if err := msg.Prepare(); err != nil {
		return nil, err
	}

	return msg.prepared, nil
}

func (msg *Message) ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error) {

	if err := msg.prepareBuffered(buffers); err != nil {
		return nil, err
	}

	return msg.prepared, nil
}

// String returns a string representation of a specified message.
func (msg *Message) String() string {
	message := Items{
		msg.beginString,
		msg.bodyLength,
		msg.msgType,
		msg.header,
	}

	if msg.body != nil {
		message = append(message, msg.body...)
	}

	if msg.trailer != nil {
		message = append(message, msg.trailer)
	}

	message = append(message, msg.checkSum)

	var items []string
	for _, item := range message {
		itemStr := item.String()
		if itemStr != "" {
			items = append(items, itemStr)
		}
	}

	return fmt.Sprintf("{%s}", strings.Join(items, ", "))
}

// Get returns an Item corresponding to the message body (identified by the item sequence number).
func (msg *Message) Get(id int) Item { return msg.body[id] }

// Set replaces an Item corresponding to the message body (identified by the item sequence number).
func (msg *Message) Set(id int, item Item) *Message { msg.body[id] = item; return msg }

// SetHeader specifies the message header.
func (msg *Message) SetHeader(header *Component) *Message { msg.setHeader(header); return msg }

// SetBody specifies the message body.
func (msg *Message) SetBody(body ...Item) *Message { msg.setBody(body); return msg }

// SetTrailer specifies the message trailer fields (except for checkSum).
func (msg *Message) SetTrailer(trailer *Component) *Message { msg.setTrailer(trailer); return msg }

func (msg *Message) setHeader(header *Component)   { msg.header = header }
func (msg *Message) setTrailer(trailer *Component) { msg.trailer = trailer }
func (msg *Message) setBody(body []Item)           { msg.body = body }

func makeTagValue(tag string, value []byte) []byte {
	return bytes.Join([][]byte{[]byte(tag), value}, []byte{61})
}
