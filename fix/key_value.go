package fix

import (
	"bytes"
	"fmt"
)

type KeyValue struct {
	Key   string
	Value Value
}

func NewKeyValue(key string, value Value) *KeyValue {
	return &KeyValue{Key: key, Value: value}
}

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

func (kv *KeyValue) Set(value Value) {
	kv.Value = value
}

func (kv *KeyValue) Load() Value {
	return kv.Value
}

func (kv *KeyValue) FromBytes(d []byte) error {
	return kv.Value.FromBytes(d)
}

func (kv *KeyValue) String() string {
	if kv.Value.IsNull() {
		return ""
	}
	return fmt.Sprintf("%s: %s", kv.Key, kv.Value)
}

type KeyValues []*KeyValue

func (kvs KeyValues) ToBytes() []byte {
	var msg [][]byte
	for _, kv := range kvs {
		if len(kv.Value.ToBytes()) > 0 {
			msg = append(msg, kv.ToBytes())
		}
	}

	return joinBody(msg...)
}
