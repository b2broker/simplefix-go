package fix

import (
	"bytes"
	"fmt"
	"reflect"
)

// UnmarshalItems parses the FIX message data stored as a byte array
// and writes it into the Items object.
func UnmarshalItems(data []byte, msg Items, strict bool) error {
	u := &unmarshaler{data: data, strict: strict}

	for _, item := range msg {
		err := u.unmarshal(u.data, item)
		if err != nil {
			return fmt.Errorf("Unmarshal items: %s", err)
		}
	}

	return nil
}

type unmarshaler struct {
	data   []byte
	strict bool
}

// scanKeyValue parses the message data related to key-value pairs
// and writes it into KeyValue objects.
func (u *unmarshaler) scanKeyValue(data []byte, el *KeyValue) error {
	q := bytes.Join([][]byte{[]byte(el.Key), {'='}}, nil)
	var keyIndex int
	if bytes.Equal(data[:len(q)], q) {
		keyIndex = 0
	} else {
		ks := bytes.Join([][]byte{Delimiter, []byte(el.Key), {'='}}, nil)
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
		return fmt.Errorf("Could not unmarshal element %s into %s: %s", el.Key, string(v), err)
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
func (u *unmarshaler) unmarshal(data []byte, fixItem Item) error {
	switch el := fixItem.(type) {
	case *KeyValue:
		return u.scanKeyValue(data, el)

	case *Group:
		noTag := el.NoTag()

		noKv := NewKeyValue(noTag, &Int{})
		err := u.unmarshal(data, noKv)
		if err != nil {
			return fmt.Errorf("Unmarshal group: %s", err)
		}

		cnt := noKv.Value.Value().(int)
		startNoTag := bytes.Index(data, append([]byte(noKv.Key), '='))
		if startNoTag == -1 {
			return nil
		}

		startFirstFieldTag := bytes.Index(data[startNoTag:], Delimiter)
		arrayString := data[startNoTag+startFirstFieldTag:]
		endFirstFieldTag := bytes.Index(arrayString, []byte{'='})

		firstTag := arrayString[:endFirstFieldTag+1]
		arrayItems := splitGroup(arrayString, firstTag)

		if len(arrayItems) == 0 {
			return fmt.Errorf("No elements found in the array")
		}

		if len(arrayItems) != cnt {
			return fmt.Errorf("Wrong items count: %d != %d", cnt, len(arrayItems))
		}

		for i := 0; i < cnt; i++ {
			entry := el.AsTemplate()

			for _, item := range entry {
				err = u.unmarshal(arrayItems[i], item)
				if err != nil {
					return fmt.Errorf("Unmarshal group item: %s", err)
				}
			}
			el.AddEntry(entry)
		}

	case *Component:
		component := el.items
		for _, item := range component {
			err := u.unmarshal(data, item)
			if err != nil {
				return fmt.Errorf("Unmarshal component: %s", err)
			}
		}

	default:
		return fmt.Errorf("Unexpected FIX item type: %s %s", reflect.TypeOf(fixItem), fixItem)
	}

	return nil
}
