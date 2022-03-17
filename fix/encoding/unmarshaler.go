package encoding

import (
	"bytes"
	"fmt"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
	"reflect"
)

type Validator interface {
	Do(msg messages.Builder) error
}

type DefaultUnmarshaller struct {
	Validator Validator
	Strict    bool
}

func NewDefaultUnmarshaller(strict bool) *DefaultUnmarshaller {
	return &DefaultUnmarshaller{Strict: strict, Validator: DefaultValidator{}}
}

func (u DefaultUnmarshaller) Unmarshal(msg messages.Builder, d []byte) error {
	err := unmarshalItems(msg.Items(), d, u.Strict)
	if err != nil {
		return err
	}

	return u.Validator.Do(msg)
}

func Unmarshal(msg messages.Builder, d []byte) error {
	u := DefaultUnmarshaller{Strict: true, Validator: DefaultValidator{}}

	return u.Unmarshal(msg, d)
}

// unmarshalItems parses the FIX message data stored as a byte array
// and writes it into the Items object.
func unmarshalItems(msg fix.Items, data []byte, strict bool) error {
	s := newState(data, strict)

	for _, item := range msg {
		err := s.unmarshal(s.data, item)
		if err != nil {
			return fmt.Errorf("could not unmarshal items: %s", err)
		}
	}

	return nil
}

type state struct {
	data   []byte
	strict bool
}

func newState(data []byte, strict bool) *state {
	return &state{data: data, strict: strict}
}

// scanKeyValue parses the message data related to key-value pairs
// and writes it into KeyValue objects.
func (s *state) scanKeyValue(data []byte, el *fix.KeyValue) error {
	q := bytes.Join([][]byte{[]byte(el.Key), {'='}}, nil)
	var keyIndex int
	if bytes.Equal(data[:len(q)], q) {
		keyIndex = 0
	} else {
		ks := bytes.Join([][]byte{fix.Delimiter, []byte(el.Key), {'='}}, nil)
		keyIndex = bytes.Index(data, ks)
		if keyIndex == -1 {
			return nil
		}
		keyIndex++ // An SOH character that is used to delimit key-value groups.
	}

	from := keyIndex + len(q)

	d := data[from:]

	end := bytes.Index(d, []byte{1})
	if end == -1 {
		end = len(d)
	}
	v := d[:end]
	err := el.FromBytes(v)
	if err != nil {
		return fmt.Errorf("could not unmarshal element %s into %s: %s", el.Key, string(v), err)
	}

	return nil
}

// splitGroup splits message parts which are recognized to be separate groups
// to create individual group items. The function distinguishes repeated parts and detects
// identical tags without repeating key-value groups.
func splitGroup(line []byte, firstTag []byte) (array [][]byte) {
	ok := true
	var index int
	for ok {
		next := bytes.Index(line[1:], firstTag)
		if next == -1 {
			index = len(line)
			ok = false
		} else {
			index = next + 1
		}
		array = append(array, line[:index])
		line = line[next+1:]
	}
	return array
}

// unmarshal traverses through a fixItem and parses its byte data,
// which is then assigned to the fixItem. A fixItem is a constructed FIX message (or its portion)
// with assigned KeyValue, Component and Group items.
func (s *state) unmarshal(data []byte, fixItem fix.Item) error {
	switch el := fixItem.(type) {
	case *fix.KeyValue:
		return s.scanKeyValue(data, el)

	case *fix.Group:
		noTag := el.NoTag()

		noKv := fix.NewKeyValue(noTag, &fix.Int{})
		err := s.unmarshal(data, noKv)
		if err != nil {
			return fmt.Errorf("could not unmarshal group: %s", err)
		}

		cnt := noKv.Value.Value().(int)
		startNoTag := bytes.Index(data, append([]byte(noKv.Key), '='))
		if startNoTag == -1 {
			return nil
		}

		startFirstFieldTag := bytes.Index(data[startNoTag:], fix.Delimiter)
		arrayString := data[startNoTag+startFirstFieldTag:]
		endFirstFieldTag := bytes.Index(arrayString, []byte{'='})

		firstTag := arrayString[:endFirstFieldTag+1]
		arrayItems := splitGroup(arrayString, firstTag)

		if len(arrayItems) == 0 {
			return fmt.Errorf("no elements found in the array")
		}

		if len(arrayItems) != cnt {
			return fmt.Errorf("wrong items count: %d != %d", cnt, len(arrayItems))
		}

		for i := 0; i < cnt; i++ {
			entry := el.AsTemplate()

			for _, item := range entry {
				err = s.unmarshal(arrayItems[i], item)
				if err != nil {
					return fmt.Errorf("could not unmarshal group item: %s", err)
				}
			}
			el.AddEntry(entry)
		}

	case *fix.Component:
		component := el.Items()
		for _, item := range component {
			err := s.unmarshal(data, item)
			if err != nil {
				return fmt.Errorf("could not unmarshal component: %s", err)
			}
		}

	default:
		return fmt.Errorf("unexpected FIX item type: %s %s", reflect.TypeOf(fixItem), fixItem)
	}

	return nil
}
