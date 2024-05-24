package gen

type SqlJoinOn struct {
	SourceTable string
	SourceAttr  string
	TargetTable string
	TargetAttr  string
}

type SqlJoinLine struct {
	Table string
	On    []SqlJoinOn
}

type SqlJoinClause struct {
	Joins []SqlJoinLine
}

type SqlSelectField struct {
	Table string
	Attr  string
	As    string
}

type SqlSelectClause struct {
	Fields []SqlSelectField
}

type SqlSelect struct {
	Table  string
	Select SqlSelectClause
	Joins  SqlJoinClause
}

type SqlGenerator interface {
	JoinOn(on SqlJoinOn) string
	JoinLine(line SqlJoinLine) string
	JoinClause(clause SqlJoinClause) string
	SelectField(field SqlSelectField) string
	SelectClause(clause SqlSelectClause) string
	Select(s SqlSelect) string
}

