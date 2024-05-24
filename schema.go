package main

import (
	"fmt"
	"strings"

	sqlparse "github.com/alicebob/sqlittle/sql"
)

type ColumnType = string

const (
	TEXT ColumnType = "text"
	INT             = "integer"
	REAL            = "real"
)

type Column struct {
	Name     string
	Type     ColumnType
	Nullable bool
}

type ForeignColumn struct {
	SourceColumn int
	TargetColumn int
}

type ForeignKey struct {
	TargetTable int
	On          []ForeignColumn
}

type Table struct {
	Name        string
	Columns     []Column
	PrimaryKey  []int
	ForeignKeys []ForeignKey
}

type Schema struct {
	Tables []Table
}

func (s Schema) FindTableIdx(name string) int {
	for i, t := range s.Tables {
		if t.Name == name {
			return i
		}
	}
	return -1
}

func (s Schema) MakeJoinClause(source, target int) string {
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

func removeCommentLines(block string) string {
	lines := strings.Split(block, "\n")
	result := ""
	for _, l := range lines {
		if strings.HasPrefix(l, "--") {
			continue
		}
		result += l + "\n"
	}
	return result
}

func parseSqlScript(script string) ([]any, error) {
	statements := []any{}

	lines := strings.Split(script, ";")
	for _, line := range lines {
		trimmed := strings.Trim(line, " \n\t")
		if trimmed == "" {
			continue
		}

		stmt, err := sqlparse.Parse(removeCommentLines(trimmed))
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			continue
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func NewSchema(schemaFile []byte) (Schema, error) {
	statements, err := parseSqlScript(string(schemaFile))
	if err != nil {
		return Schema{}, err
	}

	tables := []Table{}

	for _, s := range statements {
		switch s.(type) {
		case sqlparse.CreateTableStmt:
		default:
			continue
		}

		stmt := s.(sqlparse.CreateTableStmt)

		var pkey []int

		columns := []Column{}
		for i, col := range stmt.Columns {
			columns = append(columns, Column{
				Name:     col.Name,
				Type:     col.Type,
				Nullable: col.Null,
			})
			if col.PrimaryKey {
				pkey = []int{i}
			}
		}

		fkeys := []ForeignKey{}
		for _, intf := range stmt.Constraints {
			switch cnstr := intf.(type) {
			case sqlparse.TableForeignKey:
				targetTableIdx := -1
				for i, table := range tables {
					if table.Name == cnstr.Clause.ForeignTable {
						targetTableIdx = i
						break
					}
				}
				if targetTableIdx < 0 {
					return Schema{}, fmt.Errorf(
						"undefined target table \"%s\" of foreign key in \"%s\"",
						cnstr.Clause.ForeignTable,
						stmt.Table,
					)
				}
				targetTable := tables[targetTableIdx]

				foreignColumns := []ForeignColumn{}
				for i, col := range cnstr.Clause.ForeignColumns {
					sourceCol := cnstr.Columns[i]
					sourceColIdx := -1
					for i, col := range stmt.Columns {
						if col.Name == sourceCol {
							sourceColIdx = i
							break
						}
					}
					if sourceColIdx < 0 {
						return Schema{}, fmt.Errorf(
							"undefined source column \"%s\" of foreign key in \"%s\"",
							sourceCol, stmt.Table,
						)
					}

					targetColIdx := -1
					for i, targetCol := range targetTable.Columns {
						if targetCol.Name == col {
							targetColIdx = i
							break
						}
					}
					if targetColIdx < 0 {
						return Schema{}, fmt.Errorf(
							"undefined target column \"%s\" in target table \"%s\" of foreign key in \"%s\"",
							col,
							cnstr.Clause.ForeignTable,
							stmt.Table,
						)
					}

					foreignColumns = append(foreignColumns, ForeignColumn{
						SourceColumn: sourceColIdx,
						TargetColumn: targetColIdx,
					})
				}

				fkeys = append(fkeys, ForeignKey{
					TargetTable: targetTableIdx,
					On:          foreignColumns,
				})
			case sqlparse.TablePrimaryKey:
				pkey = []int{}
			targetColumns:
				for _, indexed := range cnstr.IndexedColumns {
					for i, col := range columns {
						if indexed.Column == col.Name {
							pkey = append(pkey, i)
							continue targetColumns
						}
					}
					return Schema{}, fmt.Errorf(
						"undefined column \"%s\" of primary key in \"%s\"",
						indexed,
						stmt.Table,
					)
				}
			}
		}

		tables = append(tables, Table{
			Name:        stmt.Table,
			Columns:     columns,
			PrimaryKey:  pkey,
			ForeignKeys: fkeys,
		})
	}

	return Schema{Tables: tables}, nil
}
