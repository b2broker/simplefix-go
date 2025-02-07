package fix

import (
	"bytes"
	"fmt"
	"strings"
)

// Item is an interface providing a method required to implement basic FIX item functionality.
type Item interface {
	ToBytes() []byte
	WriteBytes(writer *bytes.Buffer) bool
	String() string
	IsEmpty() bool
}

// Items is an array of Item elements.
type Items []Item

// ToBytes returns a byte representation of an Items array.
func (v Items) ToBytes() []byte {
	var msg [][]byte
	for _, item := range v {
		itemB := item.ToBytes()
		if itemB != nil {
			msg = append(msg, itemB)
		}
	}
	return joinBody(msg...)
}

func (v Items) IsEmpty() bool {
	for _, item := range v {
		if !item.IsEmpty() {
			return false
		}
	}
	return true
}

func (v Items) WriteBytes(writer *bytes.Buffer) bool {
	addDelimeter := false
	written := false
	for i, item := range v {
		if !item.IsEmpty() && addDelimeter {
			_ = writer.WriteByte(DelimiterChar)
			addDelimeter = false
		}
		if item.WriteBytes(writer) {
			written = true
			if i <= len(v)-1 {
				addDelimeter = true
			}
		}
	}

	return written
}

// String returns a string representation of an Items array.
func (v Items) String() string {
	var items []string
	for _, item := range v {
		itemStr := item.String()
		if itemStr != "" {
			items = append(items, itemStr)
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ", "))
}
