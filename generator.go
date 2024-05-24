package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type compositeQueryGenerator struct {
	schema Schema
}

func (g compositeQueryGenerator) getColumns(query CompositeQueryClause, table Table) []Column {
	columns := []Column{}
	for col, enabled := range query.Columns {
		if enabled {
			columns = append(columns, table.Columns[table.MustFindColumnIdx(col)])
		}
	}
	if len(columns) > 0 {
		return columns
	}

	for _, col := range table.Columns {
		columns = append(columns, col)
	}
	return columns
}

func (g compositeQueryGenerator) getGoQueryStatementName(query CompositeQuery) string {
	return fmt.Sprintf("stmt%s", query.Name)
}

// func (g compositeQueryGenerator) generateGoDataParser(query CompositeQuery) {
// 	stmtName := g.getGoQueryStatementName(query)
// }

func (g compositeQueryGenerator) generateGoDataStructureClause(query CompositeQueryClause, tableName string) string {
	properties := "{\n"

	table := g.schema.Tables[g.schema.MustFindTableIdx(tableName)]
	columns := g.getColumns(query, table)
	for _, col := range columns {
		properties += fmt.Sprintf("%s %s\n", ToUpperGoIdentifier(col.Name), SqliteToGoType(col.Type))
	}

	for childTable, childQuery := range query.With {
		properties += fmt.Sprintf(
			"%s struct %s\n",
			ToUpperGoIdentifier(childTable),
			g.generateGoDataStructureClause(childQuery, childTable),
		)
	}

	properties += "}"
	return properties
}

func (g compositeQueryGenerator) generateGoDataStructure(def CompositeQuery) string {
	return fmt.Sprintf(
		"type %sRow struct %s",
		def.Name,
		g.generateGoDataStructureClause(def.Query, def.Table),
	)
}

func (g compositeQueryGenerator) generateSelectFields(query CompositeQueryClause, tableName string) string {
	selects := ""

	table := g.schema.Tables[g.schema.MustFindTableIdx(tableName)]
	columns := g.getColumns(query, table)
	for _, col := range columns {
		selects += fmt.Sprintf(
			"    %s.%s as %s_%s,\n",
			tableName,
			col.Name,
			tableName,
			col.Name,
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

func (g compositeQueryGenerator) generateJoins(clause CompositeQueryClause, sourceTable string) string {
	joins := ""
	for targetTable, query := range clause.With {
		joins += "    " + g.generateJoinClause(
			g.schema.MustFindTableIdx(sourceTable),
			g.schema.MustFindTableIdx(targetTable),
		) + "\n"
		joins += g.generateJoins(query, targetTable)
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
		structure := g.generateGoDataStructure(def)
		err := os.WriteFile("out.go", []byte(fmt.Sprintf("package main\n\n%s", structure)), 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}
