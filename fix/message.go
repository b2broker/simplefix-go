package fix

import (
	"bytes"
	"fmt"
	"strings"
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

// BodyLength returns the BodyLength field value.
func (msg *Message) BodyLength() int {
	return msg.bodyLength.Value.Value().(int)
}

// MsgType returns the MsgType field value.
func (msg *Message) MsgType() string {
	return msg.msgType.Value.String()
}

// CheckSum returns the CheckSum field value.
func (msg *Message) CheckSum() string {
	return msg.checkSum.Value.String()
}

func (msg *Message) CalcBodyLength() int {
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()
	mt := msg.msgType.ToBytes()
	if len(bb) == 0 {
		return len(bh) + len(mt) + len(bb) + CountOfSOHSymbolsWithoutBody
	}
	return len(bh) + len(mt) + len(bb) + CountOfSOHSymbols
}

func (msg *Message) BytesWithoutChecksum() []byte {
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()

	bm := bytes.Join([][]byte{
		msg.beginString.ToBytes(),
		msg.bodyLength.ToBytes(),
		msg.msgType.ToBytes(),
		bh,
	}, Delimiter)

	if len(bb) > 0 {
		bm = bytes.Join([][]byte{bm, bb}, Delimiter)
	}

	return bm
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

// Prepared returns a byte representation of a specified message.
func (msg *Message) ToBytes() ([]byte, error) {
	err := msg.Prepare()
	if err != nil {
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
