package gen

import (
	"fmt"
	"strings"
)

type SqliteGenerator struct {}

func (g SqliteGenerator) JoinOn(on SqlJoinOn) string {
	return fmt.Sprintf(
		"%s.%s = %s.%s",
		on.SourceTable,
		on.SourceAttr,
		on.TargetTable,
		on.TargetAttr,
	)
}

func (g SqliteGenerator) JoinLine(line SqlJoinLine) string {
	clause := []string{}
	for _, on := range line.On {
		clause = append(clause, g.JoinOn(on))
	}
	return fmt.Sprintf(
		"inner join %s on %s",
		line.Table,
		strings.Join(clause, " and "),
	)
}

func (g SqliteGenerator) JoinClause(clause SqlJoinClause) string {
	result := ""
	for i, join := range clause.Joins {
		if i > 0 {
			result += "\n"
		}
		result += g.JoinLine(join)
	}
	return result
}

func (g SqliteGenerator) SelectField(field SqlSelectField) string {
	return fmt.Sprintf("%s.%s as %s", field.Table, field.Attr, field.As)
}

func (g SqliteGenerator) SelectClause(clause SqlSelectClause) string {
	if len(clause.Fields) == 0 {
		return " * "
	}
	result := ""
	for i, field := range clause.Fields {
		if i > 0 {
			result += "\n"
		}
		result += g.SelectField(field)
	}
	return result
}

func (g SqliteGenerator) Select(s SqlSelect) string {
	return fmt.Sprintf(
		"select\n%s\nfrom %s\n%s",
		g.SelectClause(s.Select),
		s.Table,
		g.JoinClause(s.Joins),
	)
}

