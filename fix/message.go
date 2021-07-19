package fix

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	// beginString, bodyLength, msgType is required fields of any message
	beginString *KeyValue
	bodyLength  *KeyValue
	msgType     *KeyValue // This tag must appear third in the list of header tags.

	// header of message expect of required fields [ bodyLength, beginString, msgType ]
	header *Component

	// body of message
	body Items

	// trailer of message expect of checkSum
	trailer *Component

	// checkSum auto generated checksum is required field
	checkSum *KeyValue

	raw []byte
}

func NewMessage(beginStringTag, bodyLengthTag, checkSumTag, msgTypeTag, beginString, msgType string) *Message {
	return &Message{
		beginString: NewKeyValue(beginStringTag, NewString(beginString)),
		bodyLength:  NewKeyValue(bodyLengthTag, &Int{}),
		msgType:     NewKeyValue(msgTypeTag, NewString(msgType)),
		checkSum:    NewKeyValue(checkSumTag, &String{}),
	}
}

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
		return fmt.Errorf("empty required field: %s (%s)", msg.beginString.Key, "BeginString")
	}
	if msg.bodyLength.Value.IsNull() {
		return fmt.Errorf("empty required field: %s (%s)", msg.bodyLength.Key, "BodyLength")
	}
	if msg.msgType.Value.IsNull() {
		return fmt.Errorf("empty required field: %s (%s)", msg.msgType.Key, "MsgType")
	}
	if msg.checkSum.Value.IsNull() {
		return fmt.Errorf("empty required field: %s (%s)", msg.checkSum.Key, "CheckSum")
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
		return fmt.Errorf("invalid body length, definde: %d, got: %d",
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
		return fmt.Errorf("invalid checksum, defined: %s, got: %s", msg.checkSum.Value, checkSum)
	}

	return nil
}

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

	if err := msg.validate(); err != nil {
		return err
	}

	return nil
}

func (msg *Message) Body() (kvs Items) {
	return msg.body
}

func (msg *Message) Header() *Component {
	return msg.header
}

func (msg *Message) Trailer() *Component {
	return msg.trailer
}

func (msg *Message) BeginString() string {
	return msg.beginString.Value.String()
}

func (msg *Message) BodyLength() int {
	return msg.bodyLength.Value.Value().(int)
}

func (msg *Message) MsgType() string {
	return msg.msgType.Value.String()
}

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

func (msg *Message) Raw() ([]byte, error) {
	if len(msg.raw) > 0 {
		return msg.raw, nil
	}

	return msg.ToBytes()
}

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
		//{'\n'},
	}, Delimiter)
	msg.raw = append(msg.raw, Delimiter...)

	if err := msg.checkRequiredFields(); err != nil {
		return nil, err
	}

	return msg.raw, nil
}

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

func (msg *Message) Get(id int) Item                        { return msg.body[id] }
func (msg *Message) Set(id int, item Item) *Message         { msg.body[id] = item; return msg }
func (msg *Message) SetRaw(raw []byte) *Message             { msg.raw = raw; return msg }
func (msg *Message) SetHeader(header *Component) *Message   { msg.setHeader(header); return msg }
func (msg *Message) SetBody(body ...Item) *Message          { msg.setBody(body); return msg }
func (msg *Message) SetTrailer(trailer *Component) *Message { msg.setTrailer(trailer); return msg }

func (msg *Message) setHeader(header *Component)   { msg.header = header }
func (msg *Message) setTrailer(trailer *Component) { msg.trailer = trailer }
func (msg *Message) setBody(body []Item)           { msg.body = body }
