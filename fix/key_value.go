package fix

import (
	"bytes"
	"fmt"
)

// KeyValue is a basic structure used for FIX message implementation.
// It is used to specify the tag and value for each field.
type KeyValue struct {
	Key   string
	Value Value
}

// NewKeyValue returns a new KeyValue object.
func NewKeyValue(key string, value Value) *KeyValue {
	return &KeyValue{Key: key, Value: value}
}

// AsTemplate returns a copy of a KeyValue object with an empty value assigned to it.
func (kv *KeyValue) AsTemplate() *KeyValue {
	switch kv.Value.(type) {
	case *String:
		return NewKeyValue(kv.Key, &String{})
	case *Int:
		return NewKeyValue(kv.Key, &Int{})
	case *Uint:
		return NewKeyValue(kv.Key, &Uint{})
	case *Time:
		return NewKeyValue(kv.Key, &Time{})
	case *Float:
		return NewKeyValue(kv.Key, &Float{})
	default:
		return NewKeyValue(kv.Key, &Raw{})
	}
}

// ToBytes returns a byte representation of a KeyValue.
func (kv *KeyValue) ToBytes() []byte {
	if kv.Value.IsNull() {
		return nil
	}

	v := kv.Value.ToBytes()
	if v == nil {
		return nil
	}

	return bytes.Join([][]byte{
		[]byte(kv.Key), v,
	}, []byte{61})
}

// Set replaces a specified value.
func (kv *KeyValue) Set(value Value) {
	kv.Value = value
}

// Load returns a specified value.
func (kv *KeyValue) Load() Value {
	return kv.Value
}

// FromBytes replaces a KeyValue object specified in the form of a byte array.
func (kv *KeyValue) FromBytes(d []byte) error {
	return kv.Value.FromBytes(d)
}

// String returns a string representation of a KeyValue object.
func (kv *KeyValue) String() string {
	if kv.Value.IsNull() {
		return ""
	}
	return fmt.Sprintf("%s: %s", kv.Key, kv.Value)
}

// KeyValues is an array of KeyValue objects.
type KeyValues []*KeyValue

// ToBytes returns a byte representation of a KeyValues array.
func (kvs KeyValues) ToBytes() []byte {
	var msg [][]byte
	for _, kv := range kvs {
		if len(kv.Value.ToBytes()) > 0 {
			msg = append(msg, kv.ToBytes())
		}
	}

	return joinBody(msg...)
}
