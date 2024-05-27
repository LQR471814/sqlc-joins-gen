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
	IsRowDef  bool
	RowDef    int
	Nullable  bool
	Array     bool
}

// a field definition in a struct, object typedef, or class
type PlFieldDef struct {
	// just for metadata usage
	TableFieldName string

	Name string
	Type PlType
}

// a struct, object typedef, or class
type PlRowDef struct {
	// just for metadata usage
	TableName string

	DefName string
	Fields  []PlFieldDef
}

// refers to a Table and a column in it
type PlScanEntry struct {
	RowDefIdx int
	FieldIdx  int
}

// refers to a collection of queries
type PlMethodDef struct {
	MethodName string
	RowDefs    []PlRowDef
	RootDef    int
	// defines the order of columns when scanning rows in
	ScanOrder []PlScanEntry
	Sql       string
}

// represents a single file in the target language
type PlScript struct {
	Methods []PlMethodDef
}

type PlScriptOutput struct {
	Path     string
	Contents []byte
}

// note: PL stands for "programming language"
// interface all programming language generators must fulfill
type PlGenerator interface {
	Script(cfg sqlc.CodegenTask, script PlScript) []PlScriptOutput
}
