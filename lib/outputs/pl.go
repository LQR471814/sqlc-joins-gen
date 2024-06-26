package outputs

// an enum of various common types in various programming languages
type PlPrimitive = int

const (
	INT PlPrimitive = iota
	FLOAT
	STRING
	BOOL
)

// the type part of a field definition
type PlType struct {
	Primitive PlPrimitive
	Nullable  bool
	Array     bool
}

// a field definition in a struct, object typedef, or class
type PlFieldDef struct {
	// just for metadata usage
	TableFieldName string

	Name     string
	Type     PlType
	IsRowDef bool
	RowDef   int
}

// a struct, object typedef, or class
type PlRowDef struct {
	// just for metadata usage
	TableName   string
	PrimaryKey  []*PlFieldDef
	Parent      *PlRowDef
	ParentField *PlFieldDef

	DefName string
	Fields  []*PlFieldDef
}

// refers to a Table and a column in it
type PlScanEntry struct {
	RowDef *PlRowDef
	Field  *PlFieldDef
}

type PlQueryArg struct {
	Name string
	Type PlType
}

// refers to a collection of queries
type PlMethodDef struct {
	MethodName string
	FirstOnly  bool
	RowDefs    []*PlRowDef
	RootDef    *PlRowDef
	// defines the order of columns when scanning rows in
	ScanOrder []PlScanEntry
	Args      []PlQueryArg
	Sql       string
}

// represents a single file in the target language
type PlScript struct {
	Methods []*PlMethodDef
}

// note: PL stands for "programming language"
// interface all programming language generators must fulfill
type PlGenerator interface {
	Generate(script PlScript) error
}
