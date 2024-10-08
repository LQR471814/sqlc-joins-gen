package inputs

import (
	"fmt"
	"strings"

	"github.com/lqr471814/sqlc-joins-gen/lib/types"

	sqlparse "github.com/alicebob/sqlittle/sql"
)

func removeSqliteCommentLines(block string) string {
	lines := strings.Split(block, "\n")
	result := ""
	for _, l := range lines {
		if strings.HasPrefix(l, "--") {
			continue
		}
		noncomment := strings.Split(l, "--")
		result += strings.TrimRight(noncomment[0], " \t") + "\n"
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
		commentless := removeSqliteCommentLines(trimmed)

		stmt, err := sqlparse.Parse(commentless)
		if err != nil {
			return nil, fmt.Errorf("sqlite statement parse failed: %v", err)
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

	var schema types.Schema
	for _, s := range statements {
		switch s.(type) {
		case sqlparse.CreateTableStmt:
		default:
			continue
		}
		stmt := s.(sqlparse.CreateTableStmt)
		table := &types.Table{Name: stmt.Table}
		for _, sqlcol := range stmt.Columns {
			table.Columns = append(table.Columns, &types.Column{
				Name:     sqlcol.Name,
				Type:     sqlcol.Type,
				Nullable: sqlcol.Null,
			})
		}
		schema.Tables = append(schema.Tables, table)
	}

	i := 0
	for _, s := range statements {
		switch s.(type) {
		case sqlparse.CreateTableStmt:
		default:
			continue
		}
		stmt := s.(sqlparse.CreateTableStmt)
		table := schema.Tables[i]

		var pkey []*types.Column
		var unique [][]*types.Column

		for i, sqlcol := range stmt.Columns {
			col := table.Columns[i]
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
					fields = append(fields, table.Columns[idx])
				}
				unique = append(unique, fields)
			case sqlparse.TableForeignKey:
				targetTableIdx := -1
				for i, table := range schema.Tables {
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
				targetTable := schema.Tables[targetTableIdx]

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
						SourceColumn: table.Columns[sourceColIdx],
						TargetColumn: targetTable.Columns[targetColIdx],
					})
				}

				fkeys = append(fkeys, types.ForeignKey{
					TargetTable: schema.Tables[targetTableIdx],
					On:          on,
				})
			case sqlparse.TablePrimaryKey:
				pkey = []*types.Column{}
			targetColumns:
				for _, indexed := range cnstr.IndexedColumns {
					for _, col := range table.Columns {
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

		schema.Tables[i].PrimaryKey = pkey
		schema.Tables[i].ForeignKeys = fkeys
		schema.Tables[i].UniqueFields = unique
		i++
	}

	return schema, nil
}
