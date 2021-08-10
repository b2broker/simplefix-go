package generator

var argTemplateFormat = `{{.Name}} {{.Type}}`

type constantsTemplate struct {
	Constants string
}

var constantsTemplateFormat = `const(
{{.Constants}}
)`

type constantTemplate struct {
	Name  string
	Value string
}

var constantTemplateFormat = `	{{.Name}} = "{{.Value}}"`

type argTemplate struct {
	Name string
	Type string
}

var componentTemplateFormat = `
type {{.Name}} struct {
	*fix.Component
}

func make{{.Name}}() *{{.Name}} {
	return &{{.Name}}{fix.NewComponent(
		{{.Fields}}
	)}
}

func New{{.Name}}({{.Args}}) *{{.Name}} {
	return make{{.Name}}(){{.Setters}}
}

{{.GetterSetters}}
`

type componentTemplate struct {
	Name          string
	Args          string
	Fields        string
	Setters       string
	GetterSetters string
}

var messageTemplateFormat = `
const MsgType{{.Name}} = "{{.MsgType}}"

type {{.Name}} struct {
	*fix.Message
}

func make{{.Name}}() *{{.Name}} {
	msg := &{{.Name}}{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgType{{.Name}}).
			SetBody(
				{{.Fields}}
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func New{{.Name}}({{.Args}}) *{{.Name}} {
	msg := make{{.Name}}(){{.Setters}}
	
	return msg
}

func Parse{{.Name}}(data []byte) (*{{.Name}}, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, FieldBeginString, beginString).
		SetBody(make{{.Name}}().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &{{.Name}}{
		Message: msg,
	}, nil
}

func ({{.LocalName}} *{{.Name}}) Header() *Header {
	header := {{.LocalName}}.Message.Header()

	return &Header{header}
}

func ({{.LocalName}} *{{.Name}}) HeaderBuilder() messages.HeaderBuilder {
	return {{.LocalName}}.Header()
}

func ({{.LocalName}} *{{.Name}}) Trailer() *Trailer {
	trailer := {{.LocalName}}.Message.Trailer()

	return &Trailer{trailer}
}

{{.GetterSetters}}
`

type messageTemplate struct {
	MsgType       string
	LocalName     string
	Name          string
	Args          string
	Fields        string
	Setters       string
	GetterSetters string
}

type defaultFlowMessage struct {
	messageTemplate
	FieldSetters string
}

var headerBuilderTemplate = `
func (Header) New() messages.HeaderBuilder {
	return makeHeader()
}
`

var trailerBuilderTemplate = `
func (Trailer) New() messages.TrailerBuilder {
	return makeTrailer()
}
`

var defaultFlowMessageTemplate = `
func ({{.Name}}) New() messages.{{.Name}}Builder {
	return make{{.Name}}()
}

func ({{.Name}}) Parse(data []byte) (messages.{{.Name}}Builder, error) {
	return Parse{{.Name}}(data)
}

{{.FieldSetters}}
`

type fieldConstructorTemplate struct {
	FieldName string
	FixType   string

	Type string
	Arg  string
}

var fieldCallConstructorTemplateFormat = `fix.NewKeyValue({{.FieldName}}, {{.FixType}}),`

type fieldGetterSetterTemplate struct {
	Index     int
	Name      string
	LocalName string
	Type      string

	ComponentName string
	ComponentType string
}

type enumVariantTemplate struct {
	Name  string
	Value string
}

var enumVariantTemplateFormat = `{{.Name}} string = "{{.Value}}"`

type enumTemplate struct {
	Name     string
	FixType  string
	Variants string
}

var enumTemplateFormat = `
// Enum type {{.Name}}
const (
 {{.Variants}}
)
`

var fieldGetterSetterTemplateFormat = `
func ({{.ComponentName}} *{{.ComponentType}}) {{.Name}}() {{.Type}} {
	kv := {{.ComponentName}}.Get({{.Index}})
	v := kv.(*fix.KeyValue).Load().Value()
	return v.({{.Type}})
}

func ({{.ComponentName}} *{{.ComponentType}}) Set{{.Name}}({{.LocalName}} {{.Type}}) *{{.ComponentType}} {
	kv := {{.ComponentName}}.Get({{.Index}}).(*fix.KeyValue)
	_ = kv.Load().Set({{.LocalName}})
	return {{.ComponentName}}
}
`

var defaultFieldSetterTemplateFormat = `
func ({{.ComponentName}} *{{.ComponentType}}) SetField{{.Name}}({{.LocalName}} {{.Type}}) messages.{{.ComponentType}}Builder {
	return {{.ComponentName}}.Set{{.Name}}({{.LocalName}})
}
`

var groupGetterSetterTemplateFormat = `
func ({{.ComponentName}} *{{.ComponentType}}) {{.Name}}() *{{.Type}} {
	group := {{.ComponentName}}.Get({{.Index}}).(*fix.Group)
	
	return &{{.Type}}{group}
}

func ({{.ComponentName}} *{{.ComponentType}}) Set{{.Name}}({{.LocalName}} *{{.Type}}) *{{.ComponentType}} {
	{{.ComponentName}}.Set({{.Index}}, {{.LocalName}}.Group)

	return {{.ComponentName}}
}
`

var componentGetterSetterTemplateFormat = `
func ({{.ComponentName}} *{{.ComponentType}}) {{.Name}}() *{{.Type}} {
	component := {{.ComponentName}}.Get({{.Index}}).(*fix.Component)
	
	return &{{.Type}}{component}
}

func ({{.ComponentName}} *{{.ComponentType}}) Set{{.Name}}({{.LocalName}} *{{.Type}}) *{{.ComponentType}} {
	{{.ComponentName}}.Set({{.Index}}, {{.LocalName}}.Component)

	return {{.ComponentName}}
}
`

type setterCallTemplate struct {
	Name  string
	Value string
}

var setterCallTemplateFormat = `Set{{.Name}}({{.Value}})`

type componentCallConstructorTemplate struct {
	Name   string
	Fields string
}

var componentCallConstructorTemplateFormat = `make{{.Name}}().Component,`

type groupCallConstructorTemplate struct {
	Name string
}

var groupCallConstructorTemplateFormat = `New{{.Name}}().Group,`

type groupConstructorTemplate struct {
	Name          string
	EntryName     string
	Fields        string
	NoTag         string
	GetterSetters string
}

var groupConstructorTemplateFormat = `
type {{.Name}} struct {
	*fix.Group
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		fix.NewGroup({{.NoTag}},
			{{.Fields}}
		),
	}
}

func (group *{{.Name}}) AddEntry(entry *{{.EntryName}}) *{{.Name}} {
	group.Group.AddEntry(entry.Items())

	return group
}
`

type fileTemplate struct {
	Data    string
	Pkg     string
	Imports string
}

var fileTemplateFormat = `
package {{.Pkg}}

import (
	{{.Imports}}
)

{{.Data}}
`
