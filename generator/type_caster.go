package generator

import (
	"fmt"
	"strings"
)

const (
	fixFloat  = "Float"
	fixInt    = "Int"
	fixRaw    = "Raw"
	fixBool   = "Bool"
	fixString = "String"
	fixTime   = "Time"
)

var allowedTypes = map[string]string{
	fixFloat:  "float64",
	fixInt:    "int",
	fixRaw:    "[]byte",
	fixBool:   "bool",
	fixString: "string",
	fixTime:   "time.Time",
}

func (g *Generator) initTypes() {
	g.typeCast = make(map[string]string, len(g.config.Types))

	for _, tp := range g.config.Types {
		if tp.CastType == "" {
			panic(fmt.Errorf("empty type attribute for type %s", tp))
		}

		if _, ok := allowedTypes[tp.CastType]; !ok {
			var types []string
			for tp := range allowedTypes {
				types = append(types, tp)
			}

			panic(fmt.Errorf(
				"unexpected type attribute %s, should be of if [%s]",
				tp, strings.Join(types, ", "),
			))
		}

		g.typeCast[tp.Name] = tp.CastType
	}
}

func (g *Generator) fixTypeToGo(t string) string {
	if tp, ok := allowedTypes[t]; ok {
		return tp
	}

	return "string"
}

func (g *Generator) typeToFix(t string) string {
	if tp, ok := g.typeCast[t]; ok {
		return tp
	}

	if _, ok := g.enums[t]; !ok {
		panic(fmt.Errorf("could not find type %s at map", t))
	}

	return fixRaw
}
