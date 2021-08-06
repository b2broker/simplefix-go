package generator

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	GroupItem     = "group"
	FieldItem     = "field"
	ComponentItem = "component"
)

// Generator generates structures and methods for FIX-messages from Doc file
type Generator struct {
	doc      *Doc
	config   *Config
	libPkg   string
	typeCast map[string]string

	enums      map[string]*Field
	fields     map[string]*Field
	components map[string]*Component
	groups     map[string]*ComponentMember
}

// NewGenerator creates new Generator
func NewGenerator(doc *Doc, config *Config, libPkg string) *Generator {
	return &Generator{
		doc:    doc,
		config: config,
		libPkg: libPkg,
		groups: make(map[string]*ComponentMember),
	}
}

func (g *Generator) checkName(name string) (err error) {
	if name == "" {
		return fmt.Errorf("empty name")
	}

	ok, err := regexp.Match("([^a-z0-9_]+|^[0-9].*)", []byte(name))
	if err != nil {
		return err
	}

	if ok {
		return fmt.Errorf("unexpected symbols in name")
	}

	return nil
}

func (g *Generator) write(path, data string) (err error) {
	output, err := os.Create(path)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", data, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("could not parse go file: %w", err)
	}

	err = (&printer.Config{Mode: printer.TabIndent, Tabwidth: 4}).Fprint(output, fset, f)
	if err != nil {
		return err
	}

	return
}

func (g *Generator) makeFile(data, pkg string) string {
	imports := []string{}

	if strings.Contains(data, "fix.") {
		// todo move
		imports = append(imports, `"github.com/b2broker/simplefix-go/fix"`)
	}
	if strings.Contains(data, "messages.") {
		// todo move
		imports = append(imports, `"github.com/b2broker/simplefix-go/session/messages"`)
	}
	if strings.Contains(data, "time.") {
		imports = append(imports, `"time"`)
	}

	return g.mustExecuteTemplate(fileTemplateFormat, fileTemplate{
		Data:    data,
		Pkg:     pkg,
		Imports: strings.Join(imports, "\n"),
	})
}

// Execute creates messages as separate files
func (g *Generator) Execute(outputDirPath string) (err error) {
	// todo split generated code into packages
	g.prepare()
	od := filepath.Clean(outputDirPath)

	dpkg := filepath.SplitList(od)[0]

	pkg := strings.Replace(dpkg, "-", "_", 10)
	if err = g.checkName(pkg); err != nil {
		return fmt.Errorf("can not use package name %s: %w", pkg, err)
	}

	pathFormat := filepath.Join(outputDirPath, "%s.generated.go")
	err = g.write(fmt.Sprintf(pathFormat, "header"), g.makeFile(g.makeHeader(), pkg))
	if err != nil {
		return err
	}

	err = g.write(fmt.Sprintf(pathFormat, "trailer"), g.makeFile(g.makeTrailer(), pkg))
	if err != nil {
		return err
	}

	err = g.write(fmt.Sprintf(pathFormat, "fields"), g.makeFile(g.makeFieldTypes(), pkg))
	if err != nil {
		return err
	}

	for _, enum := range g.enums {
		err = g.write(fmt.Sprintf(pathFormat, "enum_"+strings.ToLower(enum.Name)), g.makeFile(g.makeEnum(enum), pkg))
		if err != nil {
			return err
		}
	}

	for _, message := range g.doc.Messages {
		err = g.write(
			fmt.Sprintf(pathFormat, strings.ToLower(message.Name)),
			g.makeFile(g.makeMessage(message), pkg))
		if err != nil {
			return err
		}
	}

	for _, component := range g.components {
		err = g.write(
			fmt.Sprintf(pathFormat, strings.ToLower(component.Name)),
			g.makeFile(g.makeComponent(component, g.makeComponentTypeName(component.Name)), pkg))
		if err != nil {
			return err
		}
	}

	for _, group := range g.groups {
		err = g.write(
			fmt.Sprintf(pathFormat, strings.ToLower(g.makeGroupTypeName(group.Name))),
			g.makeFile(g.makeGroupConstructor(group), pkg))
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) prepare() {
	g.initTypes()

	g.fields = make(map[string]*Field)
	g.enums = make(map[string]*Field)
	for _, field := range g.doc.Fields {
		if len(field.Values) > 0 {
			g.enums[field.Name] = field
		} else {
			g.fields[field.Name] = field
		}
	}

	g.components = make(map[string]*Component)
	for _, component := range g.doc.Components {
		g.components[component.Name] = component
	}

	// grab groups
	for _, msg := range g.doc.Messages {
		for _, member := range msg.Members {
			g.grabGroups(member)
		}
	}
	for _, component := range g.doc.Components {
		for _, member := range component.Members {
			g.grabGroups(member)
		}
	}
	for _, member := range g.doc.Header.Members {
		g.grabGroups(member)
	}
	for _, member := range g.doc.Trailer.Members {
		g.grabGroups(member)
	}
}

func (g *Generator) validateComponent(component *ComponentMember) {
	if component.Name == "" {
		panic(fmt.Errorf("componentName could not be empty"))
	}

	if len(component.Members) == 0 {
		panic(fmt.Errorf("count of component items should be greted than 0"))
	}
}

func (g *Generator) grabGroups(component *ComponentMember) {
	if component.XMLName.Local == GroupItem {
		g.appendGroup(component)
	}
	for _, member := range component.Members {
		if member.XMLName.Local == GroupItem {
			g.validateComponent(member)

			g.appendGroup(member)
		}
		g.grabGroups(member)
	}
}

func (g *Generator) mustExecuteTemplate(format string, data interface{}) string {
	buffer := &bytes.Buffer{}
	t := template.Must(template.New("").Parse(format))
	err := t.Execute(buffer, data)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

func (g *Generator) makeStringConst(name, value string) string {
	return g.mustExecuteTemplate(constantTemplateFormat, constantTemplate{
		Name:  name,
		Value: value,
	})
}

func (g *Generator) isFieldExcluded(name string) bool {
	_, ok := ExcludedFields[name]
	return ok
}

func (g *Generator) makeArg(member *ComponentMember) string {
	tmp := argTemplate{
		Name: g.makeLocalName(member.Name),
	}

	switch member.XMLName.Local {
	case ComponentItem:
		tmp.Type = "*" + g.makeComponentTypeName(member.Name)
	case GroupItem:
		tmp.Type = "*" + g.makeGroupTypeName(member.Name)
	case FieldItem:
		tmp.Type = g.fixTypeToGo(g.makeType(member.Name))
	default:
		panic(fmt.Errorf(
			"unexpected item time, expect: %s, %s, %s, got '%s'",
			ComponentItem, GroupItem, FieldItem, member.XMLName.Local,
		))
	}

	g.makeComponentTypeName(member.Name)
	return g.mustExecuteTemplate(argTemplateFormat, tmp)
}

func (g *Generator) makeFieldTypes() string {
	goFields := make([]string, 0, len(g.doc.Fields))
	for _, field := range g.doc.Fields {
		goFields = append(goFields, g.makeStringConst("Field"+field.Name, field.Number))
	}
	return g.mustExecuteTemplate(constantsTemplateFormat, constantsTemplate{
		Constants: strings.Join(goFields, "\n"),
	})
}

func (g *Generator) validateHeader() {
	if len(g.doc.Header.Members) == 0 {
		panic(fmt.Errorf("header could not be empty"))
	}

	requiredFields := map[string]bool{}
	for rf := range RequiredHeaderFields {
		requiredFields[rf] = true
	}

	err := g.validateRequiredFields(g.doc.Header.Members, requiredFields)
	if err != nil {
		panic(fmt.Errorf("could not make header: %s", err))
	}
}

func (g *Generator) validateTrailer() {
	err := g.validateRequiredFields(g.doc.Trailer.Members, RequiredTrailerFields)
	if err != nil {
		panic(fmt.Errorf("could not make trailer: %s", err))
	}
}

func (g *Generator) validateRequiredFields(members []*ComponentMember, requiredFields map[string]bool) error {
	for _, member := range members {
		for reqField := range requiredFields {
			if member.Name == reqField {
				delete(requiredFields, reqField)
				break
			}
		}
	}

	if len(requiredFields) > 0 {
		var fields []string
		for reqField := range requiredFields {
			fields = append(fields, reqField)
		}
		return fmt.Errorf("not enougth required fields: [%s]", strings.Join(fields, ", "))
	}

	return nil
}

func (g *Generator) makeHeader() string {
	g.validateHeader()
	componentName := "Header"

	beginString := fmt.Sprintf("var beginString = \"%s.%s.%s\"", g.doc.Type, g.doc.Major, g.doc.Minor)
	header := g.makeComponent(g.doc.Header, componentName)
	fieldSetters := make([]string, len(RequiredHeaderFields))
	for fieldName := range RequiredHeaderFields {
		if g.isFieldExcluded(fieldName) {
			continue
		}

		field, ok := g.fields[fieldName]
		if !ok {
			field, ok = g.enums[fieldName]
			if !ok {
				panic(fmt.Errorf("default flow fieldName not found: %s", fieldName))
			}
		}

		fieldSetters = append(fieldSetters,
			g.mustExecuteTemplate(defaultFieldSetterTemplateFormat, fieldGetterSetterTemplate{
				Name:          field.Name,
				LocalName:     g.makeLocalName(field.Name),
				Type:          g.fixTypeToGo(g.makeType(field.Name)),
				ComponentName: g.makeLocalName(componentName),
				ComponentType: componentName,
			}),
		)
	}

	return strings.Join([]string{
		beginString,
		header,
		headerBuilderTemplate,
		strings.Join(fieldSetters, "\n")},
		"\n",
	)
}

func (g *Generator) makeTrailer() string {
	g.validateTrailer()

	return strings.Join([]string{
		g.makeComponent(g.doc.Trailer, "Trailer"),
		trailerBuilderTemplate,
	}, "\n")
}

func (g *Generator) makeComponent(component *Component, name string) string {
	goFields := make([]string, 0, len(component.Members))
	goArgs := make([]string, 0, len(component.Members))
	goSettersCalls := make([]string, 0, len(component.Members))
	goGetterSetters := make([]string, 0, len(component.Members))
	counter := 0
	for _, member := range component.Members {
		if g.isFieldExcluded(member.Name) {
			continue
		}

		goFields = append(goFields, g.makeCallConstructor(member))

		withArgs := member.Required == "Y"
		if withArgs {
			goArgs = append(goArgs, g.makeArg(member))
			goSettersCalls = append(goSettersCalls, g.makeSetterCall(member))
		}

		goGetterSetters = append(goGetterSetters, g.makeSetterGetterField(name, member, counter))

		counter++
	}

	var setters string
	if len(goSettersCalls) > 0 {
		setters = ".\n" + strings.Join(goSettersCalls, ".\n")
	}
	return g.mustExecuteTemplate(componentTemplateFormat, componentTemplate{
		Name:          name,
		Args:          strings.Join(goArgs, ", "),
		Fields:        strings.Join(goFields, "\n"),
		Setters:       setters,
		GetterSetters: strings.Join(goGetterSetters, "\n"),
	})
}

func (g *Generator) makeMessage(message *Component) string {
	goFields := make([]string, 0, len(message.Members))
	goArgs := make([]string, 0, len(message.Members))
	goSettersCalls := make([]string, 0, len(message.Members))
	goGetterSetters := make([]string, 0, len(message.Members))

	for i, member := range message.Members {
		goFields = append(goFields, g.makeCallConstructor(member))

		withArgs := member.Required == "Y"
		if withArgs {
			goArgs = append(goArgs, g.makeArg(member))
			goSettersCalls = append(goSettersCalls, g.makeSetterCall(member))
		}

		goGetterSetters = append(goGetterSetters, g.makeSetterGetterField(message.Name, member, i))
	}

	var setters string
	if len(goSettersCalls) > 0 {
		setters = ".\n" + strings.Join(goSettersCalls, ".\n")
	}

	messageData := messageTemplate{
		MsgType:       message.MsgType,
		LocalName:     g.makeLocalName(message.Name),
		Name:          message.Name,
		Args:          strings.Join(goArgs, ", "),
		Fields:        strings.Join(goFields, "\n"),
		Setters:       setters,
		GetterSetters: strings.Join(goGetterSetters, "\n"),
	}

	// if message in default
	if fields, ok := DefaultFlowFields[message.Name]; ok {
		var fieldSetters []string

		for _, fieldName := range fields {
			field, ok := g.fields[fieldName]
			if !ok {
				field, ok = g.enums[fieldName]
				if !ok {
					panic(fmt.Errorf("default flow fieldName not found: %s", fieldName))
				}
			}

			fieldSetters = append(fieldSetters,
				g.mustExecuteTemplate(defaultFieldSetterTemplateFormat, fieldGetterSetterTemplate{
					Name:          field.Name,
					LocalName:     g.makeLocalName(field.Name),
					Type:          g.fixTypeToGo(g.makeType(field.Name)),
					ComponentName: g.makeLocalName(message.Name),
					ComponentType: message.Name,
				}),
			)
		}

		return g.mustExecuteTemplate(
			messageTemplateFormat+defaultFlowMessageTemplate,
			defaultFlowMessage{
				messageTemplate: messageData,
				FieldSetters:    strings.Join(fieldSetters, "\n"),
			},
		)
	}

	return g.mustExecuteTemplate(messageTemplateFormat, messageData)
}

func (g *Generator) makeLocalName(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func (g *Generator) makeType(fieldName string) string {
	if field, ok := g.fields[fieldName]; ok {
		return g.typeToFix(field.Type)
	}

	if field, ok := g.enums[fieldName]; ok {
		return g.makeEnumName(field)
	}

	panic(fmt.Errorf("unexpected field: %s", fieldName))
}

func (g *Generator) makeEnumName(enum *Field) string {
	return fmt.Sprintf("Enum%s", enum.Name)
}

func (g *Generator) makeEnumVariant(enum string, value *Value) string {
	parts := strings.Split(value.Description, "_")
	name := ""
	// nolint
	// todo
	//if len(value.Description) > 30 {
	//	for _, part := range parts {
	//		name += string(part[0])
	//	}
	//} else {
	for _, part := range parts {
		name += string(part[0]) + strings.ToLower(part)[1:]
	}
	//}

	return g.mustExecuteTemplate(enumVariantTemplateFormat, enumVariantTemplate{
		Name: fmt.Sprintf("%s%s", enum, name),
		//Enum:  enum,
		Value: value.Enum,
	})
}

func (g *Generator) makeEnum(field *Field) string {
	name := g.makeEnumName(field)
	variants := make([]string, 0, len(field.Values))
	for _, value := range field.Values {
		variants = append(variants, g.makeEnumVariant(name, value))
	}

	return g.mustExecuteTemplate(enumTemplateFormat, enumTemplate{
		Name:     name,
		FixType:  g.fixTypeToGo(g.typeToFix(field.Name)),
		Variants: strings.Join(variants, "\n"),
	})
}

func (g *Generator) makeTypeConstructor(member *ComponentMember) string {
	if field, ok := g.fields[member.Name]; ok {
		return fmt.Sprintf("&fix.%s{}", g.typeToFix(field.Type))
	}

	if _, ok := g.enums[member.Name]; ok {
		return "&fix.String{}"
	}

	panic(fmt.Errorf("unexpected field: %s", member.Name))
}

func (g *Generator) makeFieldCallConstructor(member *ComponentMember) string {
	return g.mustExecuteTemplate(fieldCallConstructorTemplateFormat, fieldConstructorTemplate{
		FieldName: g.makeFieldName(member.Name),
		FixType:   g.makeTypeConstructor(member),
	})
}

func (g *Generator) makeFieldName(name string) string {
	return fmt.Sprintf("Field%s", name)
}

func (g *Generator) makeComponentTypeName(name string) string {
	return name
}

func (g *Generator) makeGroupTypeName(name string) string {
	return strings.Replace(name, "No", "", 1) + "Grp"
}

func (g *Generator) makeGroupEntryTypeName(name string) string {
	return strings.Replace(name, "No", "", 1) + "Entry"
}

func (g *Generator) makeSetterGetterField(parentName string, member *ComponentMember, index int) string {
	var name, tp string
	tmp := fieldGetterSetterTemplateFormat

	switch member.XMLName.Local {
	case ComponentItem:
		name = member.Name
		tp = g.makeComponentTypeName(member.Name)
		tmp = componentGetterSetterTemplateFormat

	case GroupItem:
		name = g.makeGroupTypeName(member.Name)
		tp = g.makeGroupTypeName(member.Name)
		tmp = groupGetterSetterTemplateFormat

	case FieldItem:
		name = member.Name
		tp = g.fixTypeToGo(g.makeType(member.Name))
	default:
		panic(fmt.Errorf(
			"unexpected item time, expect: %s, %s, %s, got '%s'",
			ComponentItem, GroupItem, FieldItem, member.XMLName.Local,
		))
	}

	return g.mustExecuteTemplate(tmp, fieldGetterSetterTemplate{
		Index:         index,
		Name:          name,
		LocalName:     g.makeLocalName(member.Name),
		Type:          tp,
		ComponentName: g.makeLocalName(parentName),
		ComponentType: parentName,
	})
}

func (g *Generator) makeSetterCall(member *ComponentMember) string {
	var name string
	switch member.XMLName.Local {
	case ComponentItem:
		name = member.Name
	case GroupItem:
		name = g.makeGroupTypeName(member.Name)
	case FieldItem:
		name = member.Name
	default:
		panic(fmt.Errorf(
			"unexpected item time, expect: %s, %s, %s, got '%s'",
			ComponentItem, GroupItem, FieldItem, member.XMLName.Local,
		))
	}

	return g.mustExecuteTemplate(setterCallTemplateFormat, setterCallTemplate{
		Name:  name,
		Value: g.makeLocalName(member.Name),
	})
}

func (g *Generator) makeComponentCallConstructor(member *ComponentMember) string {
	return g.mustExecuteTemplate(componentCallConstructorTemplateFormat, componentCallConstructorTemplate{
		Name: g.makeComponentTypeName(member.Name),
	})
}

func (g *Generator) makeGroupConstructor(group *ComponentMember) string {
	goGetterSetters := make([]string, 0, len(group.Members))
	goFields := make([]string, 0, len(group.Members))
	for i, member := range group.Members {
		goFields = append(goFields, g.makeCallConstructor(member))
		goGetterSetters = append(goGetterSetters, g.makeSetterGetterField(g.makeGroupEntryTypeName(group.Name), member, i))
	}

	return strings.Join([]string{
		g.mustExecuteTemplate(groupConstructorTemplateFormat, groupConstructorTemplate{
			EntryName: g.makeGroupEntryTypeName(group.Name),
			Name:      g.makeGroupTypeName(group.Name),
			NoTag:     g.makeFieldName(group.Name),
			Fields:    strings.Join(goFields, "\n"),
		}),
		g.mustExecuteTemplate(componentTemplateFormat, componentTemplate{
			Name:          g.makeGroupEntryTypeName(group.Name),
			Fields:        strings.Join(goFields, "\n"),
			GetterSetters: strings.Join(goGetterSetters, "\n"),
		}),
	}, "\n")
}

func (g *Generator) appendGroup(group *ComponentMember) {
	//	todo validate and check duplicates
	g.groups[group.Name] = group
}

func (g *Generator) makeGroupCallConstructor(member *ComponentMember) string {
	return g.mustExecuteTemplate(groupCallConstructorTemplateFormat, groupCallConstructorTemplate{
		Name: g.makeGroupTypeName(member.Name),
	})
}

func (g *Generator) makeCallConstructor(member *ComponentMember) string {
	switch member.XMLName.Local {
	case ComponentItem:
		return g.makeComponentCallConstructor(member)
	case GroupItem:
		return g.makeGroupCallConstructor(member)
	case FieldItem:
		return g.makeFieldCallConstructor(member)
	}

	panic(fmt.Errorf(
		"unexpected item time, expect: %s, %s, %s, got '%s'",
		ComponentItem, GroupItem, FieldItem, member.XMLName.Local,
	))
}
