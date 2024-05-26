package gen

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"slices"
	"sqlc-joins-gen/lib/querycfg"
	"sqlc-joins-gen/lib/schema"
	"sqlc-joins-gen/lib/sqlc"
)

type GenManager struct {
	Schema schema.Schema
}

// A "unique foreign key" is a foreign key where the collection
// of source columns is unique. (an exact match to a unique constraint,
// not a subset of a unique constriant)
func (m GenManager) isUniqueFkey(table schema.Table, fkey schema.ForeignKey) bool {
	if len(fkey.On) == len(table.PrimaryKey) {
		equal := true
		for _, col := range fkey.On {
			if slices.Index(table.PrimaryKey, col.SourceColumn) < 0 {
				equal = false
				break
			}
		}
		if equal {
			return true
		}
	}

	for _, uniqueCols := range table.UniqueFields {
		if len(fkey.On) == len(uniqueCols) {
			equal := true
			for _, col := range fkey.On {
				if slices.Index(uniqueCols, col.SourceColumn) < 0 {
					equal = false
					break
				}
			}
			if equal {
				return true
			}
		}
	}

	return false
}

func (m GenManager) getColumns(query querycfg.Query, table schema.Table) []schema.Column {
	columns := []schema.Column{}
	for col, enabled := range query.Columns {
		if enabled {
			columns = append(columns, table.Columns[table.MustFindColumnIdx(col)])
		}
	}

	if len(query.Columns) == 0 {
		return table.Columns
	}

	// don't exclude primary key columns from select
	for _, i := range table.PrimaryKey {
		col := table.Columns[i]
		if slices.Index(columns, col) >= 0 {
			continue
		}
		columns = append(columns, col)
	}

	return columns
}

func (m GenManager) getRowDef(methodName string, query querycfg.Query, tableIdx int) PlRowDef {
	table := m.Schema.Tables[tableIdx]
	def := PlRowDef{
		Name:       table.Name,
		MethodName: methodName,
	}

	columns := m.getColumns(query, table)
	for _, col := range columns {
		def.Fields = append(def.Fields, PlFieldDef{
			Name: col.Name,
			Type: PlType{
				Primitive: SqlColumnTypeToPlType(col.Type),
				Nullable:  col.Nullable,
			},
		})
	}

	for childTableName, childClause := range query.With {
		childTableIdx := m.Schema.MustFindTableIdx(childTableName)
		childTable := m.Schema.Tables[childTableIdx]

		// If the columns of a foreign key constraint are a subset of some unique constraint
		// present on the target table, the target table must only have single row attached
		// to any given current table.
		// > no more no less
		// The target table must have no less than a single row because inner join forces
		// matches on fields.
		// Meaning if there is a joined row, both the left and right must exist.

		// If the target table doesn't have any foreign keys pointing to the current table,
		// this must mean the current table must have a foreign key constraint that points
		// to the target table.
		// This means there must only exist one target table joined to any given current table.

		isUniqueFkey := false
		hasFkeyToCurrent := false
		for _, fkey := range childTable.ForeignKeys {
			if fkey.TargetTable != tableIdx {
				continue
			}
			hasFkeyToCurrent = true
			isUniqueFkey = m.isUniqueFkey(childTable, fkey)
			break
		}
		if !hasFkeyToCurrent {
			isUniqueFkey = true
		}

		def.Fields = append(def.Fields, PlFieldDef{
			Name: childTableName,
			Type: PlType{
				IsStruct: true,
				Struct:   m.getRowDef(methodName, childClause, childTableIdx),
				Array:    !isUniqueFkey,
			},
		})
	}

	return def
}

func (m GenManager) getSelectFields(query querycfg.Query, table schema.Table, out *[]SqlSelectField) {
	columns := m.getColumns(query, table)

	for _, col := range columns {
		*out = append(*out, SqlSelectField{
			Table: table.Name,
			Attr:  col.Name,
			As:    fmt.Sprintf("%s_%s", table.Name, col.Name),
		})
	}

	for tableName, childQuery := range query.With {
		table := m.Schema.MustFindTable(tableName)
		m.getSelectFields(childQuery, table, out)
	}
}

func (from GenManager) getJoinLine(sourceIdx, targetIdx int) SqlJoinLine {
	source := from.Schema.Tables[sourceIdx]
	target := from.Schema.Tables[targetIdx]

	for _, fkey := range source.ForeignKeys {
		if fkey.TargetTable == targetIdx {
			line := SqlJoinLine{Table: target.Name}
			for _, on := range fkey.On {
				line.On = append(line.On, SqlJoinOn{
					SourceTable: source.Name,
					SourceAttr:  source.Columns[on.SourceColumn].Name,
					TargetTable: target.Name,
					TargetAttr:  target.Columns[on.TargetColumn].Name,
				})
			}
			return line
		}
	}

	for _, fkey := range target.ForeignKeys {
		if fkey.TargetTable == sourceIdx {
			line := SqlJoinLine{Table: target.Name}
			for _, on := range fkey.On {
				line.On = append(line.On, SqlJoinOn{
					SourceTable: target.Name,
					SourceAttr:  target.Columns[on.SourceColumn].Name,
					TargetTable: source.Name,
					TargetAttr:  source.Columns[on.TargetColumn].Name,
				})
			}
			return line
		}
	}

	panic(fmt.Sprintf(
		"there is no possible join that can be formed from '%s' to '%s'",
		source.Name, target.Name,
	))
}

func (m GenManager) getSqlOrderBy(table string, query querycfg.Query) []SqlOrderBy {
	var orderBy []SqlOrderBy
	for col, order := range query.OrderBy {
		orderBy = append(orderBy, SqlOrderBy{
			Table:     table,
			Attr:      col,
			Ascending: order == querycfg.ASC,
		})
	}
	return orderBy
}

func (m GenManager) getJoins(query querycfg.Query, parentTableIdx int, out *[]SqlJoinLine) {
	for childTable, childQuery := range query.With {
		targetTableIdx := m.Schema.MustFindTableIdx(childTable)

		line := m.getJoinLine(parentTableIdx, targetTableIdx)

		line.Opts = SqlSelectOpts{
			Limit:   childQuery.Limit,
			Offset:  childQuery.Offset,
			Where:   childQuery.Where,
			OrderBy: m.getSqlOrderBy(childTable, childQuery),
		}

		*out = append(*out, line)

		m.getJoins(childQuery, targetTableIdx, out)
	}
}

func (m GenManager) getOrderBy(query querycfg.Query, parentTable string, out *[]SqlOrderBy) {
	// order is significant! sql sorts starting with the first orderby
	// therefore the innermost join must be sorted first
	for childTable, childQuery := range query.With {
		m.getOrderBy(childQuery, childTable, out)
	}
	table := m.Schema.MustFindTable(parentTable)
	for _, pkey := range table.PrimaryKey {
		*out = append(*out, SqlOrderBy{
			Table:     parentTable,
			Attr:      table.Columns[pkey].Name,
			Ascending: true,
		})
	}
}

func (m GenManager) getSelect(method querycfg.Method) SqlSelect {
	var fields []SqlSelectField
	tableIdx := m.Schema.MustFindTableIdx(method.Table)
	table := m.Schema.Tables[tableIdx]
	m.getSelectFields(method.Query, table, &fields)

	var joins []SqlJoinLine
	m.getJoins(method.Query, tableIdx, &joins)

	var orderBy []SqlOrderBy
	m.getOrderBy(method.Query, method.Table, &orderBy)
	orderBy = append(orderBy, m.getSqlOrderBy(method.Table, method.Query)...)

	return SqlSelect{
		Table:  method.Name,
		Select: fields,
		Joins:  joins,
		Opts: SqlSelectOpts{
			Limit:   method.Query.Limit,
			Offset:  method.Query.Offset,
			Where:   method.Query.Where,
			OrderBy: orderBy,
		},
	}
}

func (m GenManager) Generate(
	task sqlc.CodegenTask,
	sqlgen SqlGenerator,
	plgen PlGenerator,
	methods []querycfg.Method,
) error {
	script := PlScript{}
	for _, method := range methods {
		rowDef := m.getRowDef(
			method.Name,
			method.Query,
			m.Schema.MustFindTableIdx(method.Table),
		)
		rowDef.Name = method.Name + "Row"
		script.RowDefs = append(script.RowDefs, rowDef)
	}

	outputs := plgen.Script(task, script)
	for _, out := range outputs {
		interpolated := ""
		cursor := 0
		for _, location := range out.SqlLocations {
			interpolated += string(out.Contents[cursor:location.Location])

			var method querycfg.Method
			for _, m := range methods {
				if m.Name == location.MethodName {
					method = m
					break
				}
			}
			if method.Name == "" {
				return fmt.Errorf("unknown method '%s'", location.MethodName)
			}

			stmt := m.getSelect(method)
			interpolated += sqlgen.Select(stmt)
			cursor = location.Location
		}
		interpolated += string(out.Contents[cursor:len(out.Contents)])

		err := os.WriteFile(out.Path, []byte(interpolated), 0777)
		if err != nil {
			return err
		}
		err = os.Chdir(path.Dir(out.Path))
		if err != nil {
			return err
		}

		err = exec.Command("go", "fmt").Run()
		if err != nil {
			return err
		}
		withImports, err := exec.Command("goimports", out.Path).Output()
		if err != nil {
			slog.Warn(
				"failed to resolve imports with `goimports`, is it installed?",
				"target", out.Path,
			)
			continue
		}
		err = os.WriteFile(out.Path, []byte(withImports), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}