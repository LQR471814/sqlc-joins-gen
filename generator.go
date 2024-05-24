package main

import "fmt"

type JoinsGenerator struct {
	Schema Schema
	Joins  JoinsList
}

func (j JoinsGenerator) generateQuerySelect(query TableQuery, parentTable string) string {
	selects := ""

	for col, enabled := range query.Columns {
		if enabled {
			selects += fmt.Sprintf(
				"    %s.%s as %s_%s,\n",
				parentTable,
				col,
				parentTable,
				col,
			)
		}
	}

	if len(query.Columns) == 0 {
		idx := j.Schema.FindTableIdx(parentTable)
		if idx < 0 {
			panic(fmt.Sprintf("could not find table %s", parentTable))
		}

		for _, col := range j.Schema.Tables[idx].Columns {
			selects += fmt.Sprintf(
				"    %s.%s as %s_%s,\n",
				parentTable,
				col.Name,
				parentTable,
				col.Name,
			)
		}
	}

	for table, query := range query.With {
		selects += j.generateQuerySelect(query, table)
	}

	return selects
}

func (j JoinsGenerator) generateQueryJoins(query TableQuery, parentTable string) string {
	joins := ""
	for table, query := range query.With {
		joins += "    " + j.Schema.MakeJoinClause(
			j.Schema.FindTableIdx(parentTable),
			j.Schema.FindTableIdx(table),
		) + "\n"
		joins += j.generateQueryJoins(query, table)
	}
	return joins
}

func (j JoinsGenerator) generateQuerySql(def JoinQueryDef) string {
	joins := j.generateQueryJoins(def.Query, def.Table)
	selects := j.generateQuerySelect(def.Query, def.Table)

	if selects == "" {
		selects = " * "
	} else {
		selects = "\n" + selects[:len(selects)-2] + "\n"
	}

	result := fmt.Sprintf(
		"select%sfrom %s\n%s",
		selects,
		def.Table,
		joins,
	)
	return result[:len(result)-1]
}

func (j JoinsGenerator) Generate() {
	for _, def := range j.Joins {
		fmt.Println(j.generateQuerySql(def))
	}
}
