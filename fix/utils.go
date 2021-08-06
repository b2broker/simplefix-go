package fix

import (
	"bytes"
	"fmt"
)

// ValueByTag finds value by tag in raw fix-message
func ValueByTag(msg []byte, tag string) ([]byte, error) {
	start := bytes.Index(msg, bytes.Join([][]byte{{1}, []byte(tag), {61}}, nil))
	if len(msg) <= len(tag) {
		return nil, fmt.Errorf("could not found tag %s, too short message: %s", tag, msg)
	}
	if start == -1 && !bytes.Equal(bytes.Join([][]byte{[]byte(tag)}, nil), msg[:len(tag)]) {
		return nil, fmt.Errorf("tag %s not found", tag)
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
