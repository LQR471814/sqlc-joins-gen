package sqlite

import (
	"fmt"
	"sqlc-joins-gen/lib/schema"
	"strings"

	sqlparse "github.com/alicebob/sqlittle/sql"
)

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

func ParseSchema(sqliteScript []byte) (schema.Schema, error) {
	statements, err := parseSqlScript(string(sqliteScript))
	if err != nil {
		return schema.Schema{}, err
	}

	tables := []schema.Table{}

	for _, s := range statements {
		switch s.(type) {
		case sqlparse.CreateTableStmt:
		default:
			continue
		}

		stmt := s.(sqlparse.CreateTableStmt)

		var pkey []int

		columns := []schema.Column{}
		for i, col := range stmt.Columns {
			columns = append(columns, schema.Column{
				Name:     col.Name,
				Type:     col.Type,
				Nullable: col.Null,
			})
			if col.PrimaryKey {
				pkey = []int{i}
			}
		}

		fkeys := []schema.ForeignKey{}
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
					return schema.Schema{}, fmt.Errorf(
						"undefined target table \"%s\" of foreign key in \"%s\"",
						cnstr.Clause.ForeignTable,
						stmt.Table,
					)
				}
				targetTable := tables[targetTableIdx]

				foreignColumns := []schema.ForeignColumn{}
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
						return schema.Schema{}, fmt.Errorf(
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
						return schema.Schema{}, fmt.Errorf(
							"undefined target column \"%s\" in target table \"%s\" of foreign key in \"%s\"",
							col,
							cnstr.Clause.ForeignTable,
							stmt.Table,
						)
					}

					foreignColumns = append(foreignColumns, schema.ForeignColumn{
						SourceColumn: sourceColIdx,
						TargetColumn: targetColIdx,
					})
				}

				fkeys = append(fkeys, schema.ForeignKey{
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
					return schema.Schema{}, fmt.Errorf(
						"undefined column \"%s\" of primary key in \"%s\"",
						indexed,
						stmt.Table,
					)
				}
			}
		}

		tables = append(tables, schema.Table{
			Name:        stmt.Table,
			Columns:     columns,
			PrimaryKey:  pkey,
			ForeignKeys: fkeys,
		})
	}

	return schema.Schema{Tables: tables}, nil
}
