package inputs

import (
	"fmt"
	"sqlc-joins-gen/lib/types"
	"strings"

	sqlparse "github.com/alicebob/sqlittle/sql"
)

func removeSqliteCommentLines(block string) string {
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

func parseSqliteScript(script string) ([]any, error) {
	var statements []any

	lines := strings.Split(script, ";")
	for _, line := range lines {
		trimmed := strings.Trim(line, " \n\t")
		if trimmed == "" {
			continue
		}

		stmt, err := sqlparse.Parse(removeSqliteCommentLines(trimmed))
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

func ParseSqliteSchema(sqliteScript []byte) (types.Schema, error) {
	statements, err := parseSqliteScript(string(sqliteScript))
	if err != nil {
		return types.Schema{}, err
	}

	var tables []*types.Table

	for _, s := range statements {
		switch s.(type) {
		case sqlparse.CreateTableStmt:
		default:
			continue
		}

		stmt := s.(sqlparse.CreateTableStmt)

		var pkey []*types.Column
		var unique [][]*types.Column

		columns := []*types.Column{}
		for _, sqlcol := range stmt.Columns {
			col := &types.Column{
				Name:     sqlcol.Name,
				Type:     sqlcol.Type,
				Nullable: sqlcol.Null,
			}
			columns = append(columns, col)

			if sqlcol.PrimaryKey {
				pkey = []*types.Column{col}
			}
			if sqlcol.Unique {
				unique = append(unique, []*types.Column{col})
			}
		}

		var fkeys []types.ForeignKey
		for _, intf := range stmt.Constraints {
			switch cnstr := intf.(type) {
			case sqlparse.TableUnique:
				var fields []*types.Column
				for _, indexed := range cnstr.IndexedColumns {
					idx := -1
					for i, col := range stmt.Columns {
						if col.Name == indexed.Column {
							idx = i
							break
						}
					}
					if idx < 0 {
						return types.Schema{}, fmt.Errorf(
							"undefined source column \"%s\" of foreign key in \"%s\"",
							indexed.Column, stmt.Table,
						)
					}
					fields = append(fields, columns[idx])
				}
				unique = append(unique, fields)
			case sqlparse.TableForeignKey:
				targetTableIdx := -1
				for i, table := range tables {
					if table.Name == cnstr.Clause.ForeignTable {
						targetTableIdx = i
						break
					}
				}
				if targetTableIdx < 0 {
					return types.Schema{}, fmt.Errorf(
						"undefined target table \"%s\" of foreign key in \"%s\"",
						cnstr.Clause.ForeignTable,
						stmt.Table,
					)
				}
				targetTable := tables[targetTableIdx]

				on := []types.ForeignColumn{}
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
						return types.Schema{}, fmt.Errorf(
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
						return types.Schema{}, fmt.Errorf(
							"undefined target column \"%s\" in target table \"%s\" of foreign key in \"%s\"",
							col,
							cnstr.Clause.ForeignTable,
							stmt.Table,
						)
					}
					
					on = append(on, types.ForeignColumn{
						SourceColumn: columns[sourceColIdx],
						TargetColumn: targetTable.Columns[targetColIdx],
					})
				}

				fkeys = append(fkeys, types.ForeignKey{
					TargetTable: tables[targetTableIdx],
					On:          on,
				})
			case sqlparse.TablePrimaryKey:
				pkey = []*types.Column{}
			targetColumns:
				for _, indexed := range cnstr.IndexedColumns {
					for _, col := range columns {
						if indexed.Column == col.Name {
							pkey = append(pkey, col)
							continue targetColumns
						}
					}
					return types.Schema{}, fmt.Errorf(
						"undefined column \"%s\" of primary key in \"%s\"",
						indexed,
						stmt.Table,
					)
				}
			}
		}

		tables = append(tables, &types.Table{
			Name:         stmt.Table,
			Columns:      columns,
			PrimaryKey:   pkey,
			ForeignKeys:  fkeys,
			UniqueFields: unique,
		})
	}

	return types.Schema{Tables: tables}, nil
}
