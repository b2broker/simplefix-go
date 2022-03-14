package fix

import (
	"fmt"
	"strings"
)

// Item is an interface providing a method required to implement basic FIX item functionality.
type Item interface {
	ToBytes() []byte
	String() string
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
