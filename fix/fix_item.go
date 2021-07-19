package fix

import (
	"fmt"
	"strings"
)

type Item interface {
	ToBytes() []byte
	String() string
}

type Items []Item

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
