package main

import (
	"fmt"
	"strings"
)

type compositeQueryGenerator struct {
	schema Schema
}

func (g compositeQueryGenerator) getColumns(query CompositeQueryClause, tableName string) []string {
	columns := []string{}
	for col, enabled := range query.Columns {
		if enabled {
			columns = append(columns, col)
		}
	}
	if len(columns) > 0 {
		return columns
	}

	idx := g.schema.FindTableIdx(tableName)
	if idx < 0 {
		panic(fmt.Sprintf("could not find table '%s'", tableName))
	}

	table := g.schema.Tables[idx]
	for _, col := range table.Columns {
		columns = append(columns, col.Name)
	}
	return columns
}

func (g compositeQueryGenerator) generateGoDataStructureClause(query CompositeQueryClause, table string) string {
	structure := "{"

	columns := g.getColumns(query, table)
	for _, col := range columns {

	}
}

func (g compositeQueryGenerator) generateGoDataStructure(def CompositeQuery) string {
	return fmt.Sprintf(
		"type CompQuery%s struct %s",
		ToGoIdentifier(def.Table),
		g.generateGoDataStructureClause(def.Query, def.Table),
	)
}

func (g compositeQueryGenerator) generateSelectFields(query CompositeQueryClause, table string) string {
	selects := ""

	columns := g.getColumns(query, table)
	for _, col := range columns {
		selects += fmt.Sprintf(
			"    %s.%s as %s_%s,\n",
			table,
			col,
			table,
			col,
		)
	}

	for table, query := range query.With {
		selects += g.generateSelectFields(query, table)
	}

	return selects
}

func (g compositeQueryGenerator) generateJoinClause(source, target int) string {
	s := g.schema
	sourceTable := s.Tables[source]
	targetTable := s.Tables[target]

	for _, fkey := range sourceTable.ForeignKeys {
		if fkey.TargetTable == target {
			clause := []string{}
			for _, col := range fkey.On {
				clause = append(clause, fmt.Sprintf(
					"%s.%s = %s.%s",
					sourceTable.Name,
					sourceTable.Columns[col.SourceColumn].Name,
					targetTable.Name,
					targetTable.Columns[col.TargetColumn].Name,
				))
			}
			return fmt.Sprintf(
				"inner join %s on %s",
				targetTable.Name,
				strings.Join(clause, " and "),
			)
		}
	}

	for _, fkey := range targetTable.ForeignKeys {
		if fkey.TargetTable == source {
			clause := []string{}
			for _, col := range fkey.On {
				clause = append(clause, fmt.Sprintf(
					"%s.%s = %s.%s",
					targetTable.Name,
					targetTable.Columns[col.SourceColumn].Name,
					sourceTable.Name,
					sourceTable.Columns[col.TargetColumn].Name,
				))
			}
			return fmt.Sprintf(
				"inner join %s on %s",
				targetTable.Name,
				strings.Join(clause, " and "),
			)
		}
	}

	return ""
}

func (g compositeQueryGenerator) generateJoins(clause CompositeQueryClause, table string) string {
	joins := ""
	for table, query := range clause.With {
		joins += "    " + g.generateJoinClause(
			g.schema.FindTableIdx(table),
			g.schema.FindTableIdx(table),
		) + "\n"
		joins += g.generateJoins(query, table)
	}
	return joins
}

func (g compositeQueryGenerator) generateSql(def CompositeQuery) string {
	joins := g.generateJoins(def.Query, def.Table)
	selects := g.generateSelectFields(def.Query, def.Table)

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

func GenerateCompositeQuery(schema Schema, queries []CompositeQuery) {
	g := compositeQueryGenerator{schema: schema}
	for _, def := range queries {
		fmt.Println(g.generateSql(def))
	}
}
