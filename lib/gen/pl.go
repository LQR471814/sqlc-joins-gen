package gen

import (
	"fmt"
	"sqlc-joins-gen/lib/schema"
	"sqlc-joins-gen/lib/sqlc"
)

// an enum of various common types in various programming languages
type PlPrimitive = int

const (
	INT PlPrimitive = iota
	FLOAT
	STRING
	BOOL
)

// convert an sql column type into a primitive type
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

// TODO: make PlType not a recursive structure

// the type part of a field definition
type PlType struct {
	Primitive PlPrimitive
	IsStruct  bool
	Struct    int
	Nullable  bool
	Array     bool
}

// a field definition in a struct, object typedef, or class
type PlFieldDef struct {
	Name string
	Type PlType
}

// a struct, object typedef, or class
type PlRowDef struct {
	// just for metadata usage
	TableName  string
	MethodName string
	MethodRoot bool

	DefName string
	Fields  []PlFieldDef
}

// the order of fields and definitions should be
// the same as the order of SqlSelectField's
type PlScript struct {
	RowDefs []PlRowDef
}

type PlSqlLocation struct {
	MethodName string
	Location   int
}

type PlScriptOutput struct {
	Path              string
	Contents          []byte
	SqlEmbedLocations []PlSqlLocation
}

// note: PL stands for "programming language"
// interface all programming language generators must fulfill
type PlGenerator interface {
	Script(cfg sqlc.CodegenTask, script PlScript) []PlScriptOutput
}
