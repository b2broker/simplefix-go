package fix

import (
	"fmt"
	"strings"
)

type Component struct {
	items []Item
}

func NewComponent(items ...Item) *Component {
	return &Component{items: items}
}

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
	return joinBody(msg...)
}

func (c *Component) Get(id int) Item {
	return c.items[id].(*KeyValue)
}

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
