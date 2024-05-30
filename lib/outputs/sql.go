package outputs

// the part of a join after the `ON` keyword
type SqlJoinOn struct {
	SourceTable string
	SourceAttr  string
	TargetTable string
	TargetAttr  string
}

type SqlSelectOpts struct {
	Limit   int
	Offset  int
	Where   *string
	OrderBy []SqlOrderBy
}

// the full `join Table on ...` line
type SqlJoinLine struct {
	Table string
	On    []SqlJoinOn
	Opts  SqlSelectOpts
}

// the part of a select clause that limits what columns are returned
// looks something like `Table.column as table_column`
type SqlSelectField struct {
	Table string
	Attr  string
	As    string
}

type SqlOrderBy struct {
	Table string
	Attr  string
	// true = asc, false = dsc
	Ascending bool
}

// a single select statement `select ... from Table join ForeignTable on ...`
type SqlSelect struct {
	Table     string
	FirstOnly bool
	Select    []SqlSelectField
	Joins     []SqlJoinLine
	Opts      SqlSelectOpts
}

// the interface all sql generators must fulfill
type SqlGenerator interface {
	Select(s SqlSelect) string
}
