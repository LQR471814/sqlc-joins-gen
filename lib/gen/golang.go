package gen

import (
	"fmt"
	"sqlc-joins-gen/lib/utils"
	"strings"
)

func lowerGoIdentifier(name string) string {
	result := ""

	for i, c := range name {
		if i == 0 && c >= '0' && c <= '9' {
			result += "T"
		}
		if c == '-' || c == '_' {
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += strings.ToLower(string(c))
			continue
		}
		panic(fmt.Sprintf(
			"got invalid character in name to be converted to an identifier '%c', '%s'",
			c,
			name,
		))
	}

	return result
}

func upperGoIdentifier(name string) string {
	return utils.Capitalize(lowerGoIdentifier(name))
}

type GolangGenerator struct{}

func (g GolangGenerator) Type(t PlType) string {
	if t.IsStruct {
		return g.StructDef(t.Struct)
	}
	switch t.Primitive {
	case INT:
		return "int"
	case FLOAT:
		return "float64"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	}
	panic(fmt.Sprintf("unknown primitive type '%d'", t.Primitive))
}

func (g GolangGenerator) FieldDef(def PlFieldDef) string {
	return fmt.Sprintf(
		"%s %s",
		def.Name, g.Type(def.Type),
	)
}

func (g GolangGenerator) FieldClause(clause PlFieldClause) string {
	result := "{\n"
	for _, def := range clause.Fields {
		result += g.FieldDef(def) + "\n"
	}
	return result + "}"
}

func (g GolangGenerator) StructDef(def PlStructDef) string {
	return fmt.Sprintf(
		"type %s struct %s",
		upperGoIdentifier(def.Name),
		g.FieldClause(def.Fields),
	)
}
