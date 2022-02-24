package fix

import (
	"bytes"
	"fmt"
	"strconv"
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

	raw []byte
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

// NewMessageFromBytes creates a new empty message.
func NewMessageFromBytes(beginStringTag, bodyLengthTag, checkSumTag, msgTypeTag string) *Message {
	return &Message{
		beginString: NewKeyValue(beginStringTag, &String{}),
		bodyLength:  NewKeyValue(bodyLengthTag, &Int{}),
		msgType:     NewKeyValue(msgTypeTag, &String{}),
		checkSum:    NewKeyValue(checkSumTag, &String{}),
	}
}

func (msg *Message) checkRequiredFields() error {
	if msg.beginString.Value.IsNull() {
		return fmt.Errorf("The required field value is empty: %s (%s)", msg.beginString.Key, "BeginString")
	}
	if msg.bodyLength.Value.IsNull() {
		return fmt.Errorf("The required field value is empty: %s (%s)", msg.bodyLength.Key, "BodyLength")
	}
	if msg.msgType.Value.IsNull() {
		return fmt.Errorf("The required field value is empty: %s (%s)", msg.msgType.Key, "MsgType")
	}
	if msg.checkSum.Value.IsNull() {
		return fmt.Errorf("The required field value is empty: %s (%s)", msg.checkSum.Key, "CheckSum")
	}

	return nil
}

func (msg *Message) validate() error {
	if err := msg.checkRequiredFields(); err != nil {
		return err
	}

	mt := msg.msgType.ToBytes()
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()
	bodyLength := msg.calcBodyLength(bh, bb, mt)
	if bodyLength != msg.bodyLength.Value.Value().(int) {
		return fmt.Errorf("An invalid body length; specified: %d, required: %d",
			msg.bodyLength.Value.Value().(int),
			bodyLength,
		)
	}

	byteMsg := bytes.Join([][]byte{
		msg.beginString.ToBytes(),
		msg.bodyLength.ToBytes(),
		msg.msgType.ToBytes(),
		bh,
	}, Delimiter)

	if len(bb) > 0 {
		byteMsg = bytes.Join([][]byte{byteMsg, bb}, Delimiter)
	}

	checkSum := string(calcCheckSum(byteMsg))

	if checkSum != msg.checkSum.Value.String() {
		return fmt.Errorf("An invalid checksum; specified: %s, required: %s", msg.checkSum.Value, checkSum)
	}

	return nil
}

// Unmarshal parses byte data and inserts it into a specified Message object.
func (msg *Message) Unmarshal(data []byte) error {
	message := Items{
		msg.beginString,
		msg.bodyLength,
		msg.msgType,
		msg.header,
	}

	message = append(message, msg.body...)
	message = append(message, msg.trailer, msg.checkSum)

	err := UnmarshalItems(data, message, false)
	if err != nil {
		return err
	}

	return msg.validate()
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
func (msg *Message) BeginString() string {
	return msg.beginString.Value.String()
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

func (msg *Message) calcBodyLength(header, body, msgType []byte) int {
	count := len(header) + len(msgType) + len(body) + CountOfSOHSymbols
	if len(body) == 0 {
		return count - 1
	}

	return count
}

// Raw returns message data in the form of a byte array.
func (msg *Message) Raw() ([]byte, error) {
	if len(msg.raw) > 0 {
		return msg.raw, nil
	}

	return msg.ToBytes()
}

// ToBytes returns a byte representation of a specified message.
func (msg *Message) ToBytes() ([]byte, error) {
	bh := msg.header.ToBytes()
	bb := msg.body.ToBytes()
	mt := msg.msgType.ToBytes()
	msg.bodyLength.Value = NewString(strconv.Itoa(msg.calcBodyLength(bh, bb, mt)))

	byteMsg := bytes.Join([][]byte{
		msg.beginString.ToBytes(),
		msg.bodyLength.ToBytes(),
		msg.msgType.ToBytes(),
		bh,
	}, Delimiter)

	if len(bb) > 0 {
		byteMsg = bytes.Join([][]byte{byteMsg, bb}, Delimiter)
	}

	checkSum := calcCheckSum(byteMsg)
	err := msg.checkSum.Value.Set(string(checkSum))
	if err != nil {
		return nil, err
	}

	msg.raw = bytes.Join([][]byte{
		byteMsg,
		makeTagValue(msg.checkSum.Key, checkSum),
	}, Delimiter)
	msg.raw = append(msg.raw, Delimiter...)

	if err := msg.checkRequiredFields(); err != nil {
		return nil, err
	}

	return msg.raw, nil
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
