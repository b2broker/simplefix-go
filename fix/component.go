package fix

import (
	"bytes"
	"fmt"
	"strings"
)

// Component is an array of various FIX entities.
// It may contain a KeyValue, a Group and another Component.
type Component struct {
	items []Item
}

// NewComponent is used to create a new Component instance.
func NewComponent(items ...Item) *Component {
	return &Component{items: items}
}

// Items returns Component items.
func (c *Component) Items() Items {
	return c.items
}

// AsComponent returns a specified component.
// This is required for structures that integrate this component.
func (c *Component) AsComponent() *Component {
	return c
}

// AsTemplate returns a new structure with the same set of items
// as in a specified component (these items are assigned empty values).
func (c *Component) AsTemplate() Items {
	tmp := make([]Item, len(c.items))

	for i, item := range c.items {
		switch value := item.(type) {
		case *KeyValue:
			tmp[i] = value.AsTemplate()

		case *Group:
			tmp[i] = NewGroup(value.NoTag(), value.AsTemplate()...)

		case *Component:
			tmp[i] = NewComponent(value.AsTemplate()...)
		}
	}

	return tmp
}

// ToBytes returns a representation of a message which is native to FIX.
func (c *Component) ToBytes() []byte {
	var msg [][]byte
	for _, item := range c.items {
		itemB := item.ToBytes()
		if itemB != nil {
			msg = append(msg, itemB)
		}
	}

	if len(msg) == 0 {
		return nil
	}

	return joinBody(msg...)
}
func (c *Component) IsEmpty() bool {
	for _, item := range c.items {
		if !item.IsEmpty() {
			return false
		}
	}
	return true
}
func (c *Component) WriteBytes(writer *bytes.Buffer) bool {
	addDelimeter := false
	written := false
	for i, item := range c.items {
		if !item.IsEmpty() && addDelimeter {
			_ = writer.WriteByte(DelimiterChar)
			addDelimeter = false
		}
		if item.WriteBytes(writer) {
			written = true
			if i <= len(c.items)-1 {
				addDelimeter = true
			}
		}
	}

	return written
}

// Get returns a specific component item identified by its sequence number.
// Such item may be a *KeyValue, *Component or *Group.
func (c *Component) Get(id int) Item {
	return c.items[id]
}

// Set replaces a component item identified by its sequence number.
func (c *Component) Set(id int, v Item) {
	c.items[id] = v
}

// SetField is used to define an internal field for any item.
func (c *Component) SetField(id int, v Item) {
	c.items[id] = v
}

// SetGroup is used to define an internal group for an item.
func (c *Component) SetGroup(id int, v *Group) {
	c.items[id] = v
}

// SetComponent is used to define an internal component for an item.
func (c *Component) SetComponent(id int, v *Component) {
	c.items[id] = v
}

// String returns a string representation of a component.
func (c *Component) String() string {
	var items []string
	for _, item := range c.items {
		itemStr := item.String()
		if itemStr != "" {
			items = append(items, itemStr)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}
