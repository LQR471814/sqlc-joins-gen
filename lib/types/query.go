package types

type Method struct {
	Name   string
	Return Return
	Table  *Table
	Query  Query
}

type QueryColumn struct {
	Column  *Column
	Enabled bool
}

type QueryOrderBy struct {
	Column  *Column
	OrderBy OrderBy
}

type QueryWith struct {
	Table *Table
	Query Query
}

type Query struct {
	Columns []QueryColumn
	With    []QueryWith
	OrderBy []QueryOrderBy
	Where   *string
	Limit   int
	Offset  int
}

type Return = string

const (
	FIRST Return = "first"
	MANY         = "many"
)

type OrderBy = string

const (
	ASC OrderBy = "asc"
	DSC OrderBy = "dsc"
)
