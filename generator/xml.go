package generator

import "encoding/xml"

// Doc is a structure identifying the components and fields required for a custom FIX protocol implementation.
type Doc struct {
	Type        string `xml:"type,attr"`
	Major       string `xml:"major,attr"`
	Minor       string `xml:"minor,attr"`
	ServicePack int    `xml:"servicepack,attr"`

	Header     *Component   `xml:"header"`
	Trailer    *Component   `xml:"trailer"`
	Messages   []*Component `xml:"messages>message"`
	Components []*Component `xml:"components>component"`
	Fields     []*Field     `xml:"fields>field"`
}

// Component is a structure identifying the set of basic elements required for FIX messages,
// such as key-value groups or basic components.
type Component struct {
	Name    string `xml:"name,attr"`
	MsgCat  string `xml:"msgcat,attr"`
	MsgType string `xml:"msgtype,attr"`

	Members []*ComponentMember `xml:",any"`
}

// Field is a structure used to implement key-value groups.
type Field struct {
	Number string   `xml:"number,attr"`
	Name   string   `xml:"name,attr"`
	Type   string   `xml:"type,attr"`
	Values []*Value `xml:"value"`
}

// Value is a structure used to support enumeration-like FIX values.
type Value struct {
	Enum        string `xml:"enum,attr"`
	Description string `xml:"description,attr"`
}

type ComponentMember struct {
	XMLName  xml.Name
	Name     string `xml:"name,attr"`
	Required string `xml:"required,attr"`

	Members []*ComponentMember `xml:",any"`
}

type Config struct {
	Types []*Type `xml:"types>type"`
}

type Type struct {
	XMLName xml.Name

	Name     string `xml:"name,attr"`
	CastType string `xml:"cast,attr"`
}
