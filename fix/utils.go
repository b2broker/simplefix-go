package fix

import (
	"bytes"
	"fmt"
)

// ValueByTag locates a value by its tag in a FIX message stored as a byte array.
func ValueByTag(msg []byte, tag string) ([]byte, error) {
	start := bytes.Index(msg, bytes.Join([][]byte{{1}, []byte(tag), {61}}, nil))
	if len(msg) <= len(tag) {
		return nil, fmt.Errorf("could not find the tag: %s, the message is too short: %s", tag, msg)
	}
	if start == -1 && !bytes.Equal(bytes.Join([][]byte{[]byte(tag)}, nil), msg[:len(tag)]) {
		return nil, fmt.Errorf("the tag is not found: %s", tag)
	}
	start += len(tag) + 2
	end := bytes.Index(msg[start:], []byte{1})
	if end == -1 {
		end = len(msg)
	} else {
		end += start
	}
	return msg[start:end], nil
}
