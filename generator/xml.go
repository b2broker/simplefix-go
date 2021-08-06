package generator

import "encoding/xml"

// Doc is a structure with components and fields of special FIX-protocol
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

// Component is a minimal part of FIX-message like key-value group or component
type Component struct {
	Name    string `xml:"name,attr"`
	MsgCat  string `xml:"msgcat,attr"`
	MsgType string `xml:"msgtype,attr"`

	Members []*ComponentMember `xml:",any"`
}

// Field is a key-value component
type Field struct {
	Number string   `xml:"number,attr"`
	Name   string   `xml:"name,attr"`
	Type   string   `xml:"type,attr"`
	Values []*Value `xml:"value"`
}

// Value is a enum-variant
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
