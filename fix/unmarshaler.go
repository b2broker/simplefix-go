package fix

import (
	"bytes"
	"fmt"
	"reflect"
)

// UnmarshalItems reads data from byte FIX-message and write it in Items
func UnmarshalItems(data []byte, msg Items, strict bool) error {
	u := &unmarshaler{data: data, strict: strict}

	for _, item := range msg {
		err := u.unmarshal(u.data, item)
		if err != nil {
			return fmt.Errorf("unmarshal items: %s", err)
		}
	}

	return nil
}

type unmarshaler struct {
	data   []byte
	strict bool
}

// scanKeyValue reads data from message part to special KeyValue
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
		keyIndex++ // SOH
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
		return fmt.Errorf("could not unmarshal el %s into %s: %s", el.Key, string(v), err)
	}

	return nil
}

// splitGroup splits part of message which we think are group to group items
// it is separates repeated parts and allows to looking for the same tags without repeat.
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

// unmarshal is traversal of fixItem and scan bytes from data into fixItem
// fixItem is a prepared fix message (or part of message) with KeyValue, Component and Group items
func (u *unmarshaler) unmarshal(data []byte, fixItem Item) error {
	switch el := fixItem.(type) {
	case *KeyValue:
		return u.scanKeyValue(data, el)

	case *Group:
		noTag := el.NoTag()

		noKv := NewKeyValue(noTag, &Int{})
		err := u.unmarshal(data, noKv)
		if err != nil {
			return fmt.Errorf("unmarshal group: %s", err)
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
			return fmt.Errorf("no elements")
		}

		if len(arrayItems) != cnt {
			return fmt.Errorf("wront count of items: %d != %d", cnt, len(arrayItems))
		}

		for i := 0; i < cnt; i++ {
			entry := el.AsTemplate()

			for _, item := range entry {
				err = u.unmarshal(arrayItems[i], item)
				if err != nil {
					return fmt.Errorf("unmarshal group item: %s", err)
				}
			}
			el.AddEntry(entry)
		}

	case *Component:
		component := el.items
		for _, item := range component {
			err := u.unmarshal(data, item)
			if err != nil {
				return fmt.Errorf("unmarshal component: %s", err)
			}
		}

	default:
		return fmt.Errorf("unexpected type of fix item: %s %s", reflect.TypeOf(fixItem), fixItem)
	}

	return nil
}
