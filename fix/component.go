package fix

import (
	"fmt"
	"strings"
)

// Component is a batch of different FIX-items
// it may contain KeyValue, Group and another Component
type Component struct {
	items []Item
}

// NewComponent creates new Component
func NewComponent(items ...Item) *Component {
	return &Component{items: items}
}

// Items returns items of Component
func (c *Component) Items() Items {
	return c.items
}

// AsComponent returns itself
// need fot structures which aggregate this Component
func (c *Component) AsComponent() *Component {
	return c
}

// AsTemplate returns new structure with same items with empty values
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

// ToBytes returns FIX-representation message
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

// Get returns item of component by sequence number
// it may be *KeyValue, *Component or *Group
func (c *Component) Get(id int) Item {
	return c.items[id]
}

// Set replace item of component by sequence number
func (c *Component) Set(id int, v Item) {
	c.items[id] = v
}

// SetGroup sets internal field (any) item
func (c *Component) SetField(id int, v Item) {
	c.items[id] = v
}

// SetGroup sets internal group item
func (c *Component) SetGroup(id int, v *Group) {
	c.items[id] = v
}

// SetComponent sets internal component item
func (c *Component) SetComponent(id int, v *Component) {
	c.items[id] = v
}

// String converts Component to string
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
