package outputs

import (
	"fmt"
	"strings"
)

type SqliteGenerator struct{}

func (g SqliteGenerator) joinOn(on SqlJoinOn) string {
	return fmt.Sprintf(
		`"%s"."%s" = "%s"."%s"`,
		on.SourceTable,
		on.SourceAttr,
		on.TargetTable,
		on.TargetAttr,
	)
}

func (g SqliteGenerator) joinLine(line SqlJoinLine) string {
	var clause []string
	for _, on := range line.On {
		clause = append(clause, g.joinOn(on))
	}
	onClause := strings.Join(clause, " and ")

	if line.Opts.Limit > 0 || line.Opts.Offset > 0 ||
		(line.Opts.Where != nil && *line.Opts.Where != "") || line.Opts.OrderBy != nil {
		subselect := g.selectSubquery(line.Table, line.Opts)
		return fmt.Sprintf(
			`inner join (%s) as "%s" on %s`,
			subselect, line.Table, onClause,
		)
	}

	return fmt.Sprintf(
		`inner join "%s" on %s`,
		line.Table, onClause,
	)
}

func (g SqliteGenerator) joinLineList(lines []SqlJoinLine) string {
	result := ""
	for i, join := range lines {
		if i > 0 {
			result += "\n"
		}
		result += g.joinLine(join)
	}
	return result
}

func (g SqliteGenerator) selectField(field SqlSelectField) string {
	return fmt.Sprintf(`"%s"."%s" as "%s"`, field.Table, field.Attr, field.As)
}

func (g SqliteGenerator) selectFieldList(fields []SqlSelectField) string {
	if len(fields) == 0 {
		return " * "
	}
	result := ""
	for i, field := range fields {
		if i > 0 {
			result += ",\n"
		}
		result += g.selectField(field)
	}
	return result
}

func (g SqliteGenerator) orderByList(list []SqlOrderBy) string {
	if len(list) == 0 {
		return ""
	}
	result := "order by"
	for _, orderBy := range list {
		keyword := "asc"
		if !orderBy.Ascending {
			keyword = "dsc"
		}
		result += fmt.Sprintf(` "%s"."%s" %s,`, orderBy.Table, orderBy.Attr, keyword)
	}
	return result[:len(result)-1]
}

func (g SqliteGenerator) where(query string) string {
	if query == "" {
		return ""
	}
	return fmt.Sprintf("where %s", query)
}

func (g SqliteGenerator) limitAndOffset(limit, offset int) string {
	result := ""
	if limit > 0 {
		result += fmt.Sprintf("limit %d", limit)
	}
	if offset > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("offset %d", offset)
	}
	return result
}

func (g SqliteGenerator) selectSubquery(table string, opts SqlSelectOpts) string {
	where := g.where(*opts.Where)
	orderBy := g.orderByList(opts.OrderBy)
	limitAndOffset := g.limitAndOffset(opts.Limit, opts.Offset)
	res := fmt.Sprintf(`select * from "%s"`, table)
	if where != "" {
		res += " " + where
	}
	if orderBy != "" {
		res += " " + orderBy
	}
	if limitAndOffset != "" {
		res += " " + limitAndOffset
	}
	return res
}

func (g SqliteGenerator) Select(s SqlSelect) string {
	if s.FirstOnly {
		tableQuery := g.selectSubquery(s.Table, SqlSelectOpts{
			Limit:   1,
			Offset:  s.Opts.Offset,
			Where:   s.Opts.Where,
			OrderBy: s.Opts.OrderBy,
		})
		res := fmt.Sprintf(
			"select\n%s\nfrom (%s) %s",
			g.selectFieldList(s.Select),
			tableQuery,
			g.joinLineList(s.Joins),
		)
		if s.Opts.Limit > 0 {
			res += fmt.Sprintf(" limit %d", s.Opts.Limit)
		}
		return res
	}
	return fmt.Sprintf(
		"select\n%s\nfrom \"%s\"\n%s\n%s\n%s\n%s",
		g.selectFieldList(s.Select),
		s.Table,
		g.joinLineList(s.Joins),
		g.where(*s.Opts.Where),
		g.orderByList(s.Opts.OrderBy),
		g.limitAndOffset(s.Opts.Limit, s.Opts.Offset),
	)
}
