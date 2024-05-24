package gen

import (
	"fmt"
	"sqlc-joins-gen/lib/schema"
)

type PlPrimitive = int

const (
	INT PlPrimitive = iota
	FLOAT
	STRING
	BOOL
)

type PlType struct {
	Primitive PlPrimitive
	IsStruct  bool
	Struct    PlStructDef
}

type PlFieldDef struct {
	Name     string
	Type     PlType
	Nullable bool
}

type PlFieldClause struct {
	Fields []PlFieldDef
}

type PlStructDef struct {
	Name   string
	Fields PlFieldClause
}

type PlGenerator interface {
	Type(t PlType) string
	FieldDef(field PlFieldDef) string
	FieldClause(clause PlFieldClause) string
	StructDef(def PlStructDef) string
}

func SqlColumnTypeToPlType(t schema.ColumnType) PlPrimitive {
	switch t {
	case schema.INT:
		return INT
	case schema.TEXT:
		return STRING
	case schema.REAL:
		return FLOAT
	}
	panic(fmt.Sprintf("unknown column type '%s'", t))
}
