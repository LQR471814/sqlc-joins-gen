package transform

import (
	"fmt"
	"regexp"
	"slices"
	"sqlc-joins-gen/lib/outputs"
	"sqlc-joins-gen/lib/types"
	"sqlc-joins-gen/lib/utils"
)

type FromSchema struct {
	Schema types.Schema
}

// A "unique foreign key" is a foreign key where the collection
// of source columns is unique. (an exact match to a unique constraint,
// not a subset of a unique constriant)
func (s FromSchema) isUniqueFkey(table *types.Table, fkey types.ForeignKey) bool {
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

func (s FromSchema) getColumns(query types.Query, table *types.Table) []*types.Column {
	columns := []*types.Column{}
	for _, col := range query.Columns {
		if col.Enabled {
			columns = append(columns, col.Column)
		}
	}

	if len(query.Columns) == 0 {
		return table.Columns
	}

	// don't exclude primary key columns from select
	for _, col := range table.PrimaryKey {
		if slices.Index(columns, col) >= 0 {
			continue
		}
		columns = append(columns, col)
	}

	return columns
}

var queryArgRegex = regexp.MustCompile(`\$(\w+):(\w+)(\?)?(\[\])?`)

func (s FromSchema) getQueryArgs(query types.Query, out *[]outputs.PlQueryArg) error {
	var err error

	*query.Where = utils.ReplaceAllStringSubmatchFunc(queryArgRegex, *query.Where, func(m []string) string {
		name := m[1]
		typestr := m[2]
		nullable := m[3] != ""
		isArray := m[4] != ""

		var primitive outputs.PlPrimitive
		switch typestr {
		case "int":
			primitive = outputs.INT
		case "str":
			primitive = outputs.STRING
		case "float":
			primitive = outputs.FLOAT
		case "bool":
			primitive = outputs.BOOL
		default:
			err = fmt.Errorf(
				"unknown type in arg definition '%s' in '%s'",
				m[2], m[0],
			)
			return "ERROR"
		}

		*out = append(*out, outputs.PlQueryArg{
			Name: name,
			Type: outputs.PlType{
				Primitive: primitive,
				Nullable:  nullable,
				Array:     isArray,
			},
		})

		return fmt.Sprintf("$%d", len(*out)-1)
	})
	if err != nil {
		return err
	}

	for _, with := range query.With {
		err = s.getQueryArgs(with.Query, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s FromSchema) getRowDefs(method types.Method, out *[]*outputs.PlRowDef) {
	queue := []struct {
		rowDefName  string
		parent      *outputs.PlRowDef
		parentField *outputs.PlFieldDef
		table       *types.Table
		query       types.Query
	}{
		{
			rowDefName: method.Name,
			table:      method.Table,
			query:      method.Query,
		},
	}

	rowDefIdxOffset := len(*out)
	for len(queue) > 0 {
		current := queue[0]
		if len(queue) > 0 {
			queue = queue[1:]
		}

		def := &outputs.PlRowDef{
			TableName:   current.table.Name,
			DefName:     current.rowDefName,
			Parent:      current.parent,
			ParentField: current.parentField,
		}

		columns := s.getColumns(current.query, current.table)
		for _, col := range columns {
			field := &outputs.PlFieldDef{
				TableFieldName: col.Name,
				Name:           col.Name,
				Type: outputs.PlType{
					Primitive: SqlColumnTypeToPlType(col.Type),
					Nullable:  col.Nullable,
				},
			}

			for _, pkey := range current.table.PrimaryKey {
				if pkey == col {
					def.PrimaryKey = append(def.PrimaryKey, field)
					break
				}
			}
			def.Fields = append(def.Fields, field)
		}

		i := 0
		for _, with := range current.query.With {
			childTable := with.Table

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
				if fkey.TargetTable != current.table {
					continue
				}
				hasFkeyToCurrent = true
				isUniqueFkey = s.isUniqueFkey(childTable, fkey)
				break
			}
			if !hasFkeyToCurrent {
				isUniqueFkey = true
			}

			rowDefIdxOffset++
			field := &outputs.PlFieldDef{
				Name:     childTable.Name,
				Type:     outputs.PlType{Array: !isUniqueFkey},
				IsRowDef: true,
				RowDef:   rowDefIdxOffset,
			}

			queue = append(queue, struct {
				rowDefName  string
				parent      *outputs.PlRowDef
				parentField *outputs.PlFieldDef
				table       *types.Table
				query       types.Query
			}{
				rowDefName:  current.rowDefName + fmt.Sprint(i),
				parent:      def,
				parentField: field,
				table:       childTable,
				query:       with.Query,
			})

			def.Fields = append(def.Fields, field)
			i++
		}

		*out = append(*out, def)
	}
}

func (s FromSchema) GetMethodDef(method types.Method) (outputs.PlMethodDef, error) {
	out := outputs.PlMethodDef{
		MethodName: method.Name,
		FirstOnly:  method.Return == types.FIRST,
	}
	err := s.getQueryArgs(method.Query, &out.Args)
	if err != nil {
		return outputs.PlMethodDef{}, err
	}
	s.getRowDefs(method, &out.RowDefs)
	return out, nil
}

func (s FromSchema) getSelectFields(query types.Query, table *types.Table, out *[]outputs.SqlSelectField) {
	columns := s.getColumns(query, table)

	for _, col := range columns {
		*out = append(*out, outputs.SqlSelectField{
			Table: table.Name,
			Attr:  col.Name,
			As:    fmt.Sprintf("%s_%s", table.Name, col.Name),
		})
	}

	for _, with := range query.With {
		s.getSelectFields(with.Query, with.Table, out)
	}
}

func (s FromSchema) getJoinLine(source, target *types.Table) (outputs.SqlJoinLine, error) {
	for _, fkey := range source.ForeignKeys {
		if fkey.TargetTable == target {
			line := outputs.SqlJoinLine{Table: target.Name}
			for _, on := range fkey.On {
				line.On = append(line.On, outputs.SqlJoinOn{
					SourceTable: source.Name,
					SourceAttr:  on.SourceColumn.Name,
					TargetTable: target.Name,
					TargetAttr:  on.TargetColumn.Name,
				})
			}
			return line, nil
		}
	}

	for _, fkey := range target.ForeignKeys {
		if fkey.TargetTable == source {
			line := outputs.SqlJoinLine{Table: target.Name}
			for _, on := range fkey.On {
				line.On = append(line.On, outputs.SqlJoinOn{
					SourceTable: target.Name,
					SourceAttr:  on.SourceColumn.Name,
					TargetTable: source.Name,
					TargetAttr:  on.TargetColumn.Name,
				})
			}
			return line, nil
		}
	}

	return outputs.SqlJoinLine{}, fmt.Errorf(
		"there is no possible join that can be formed from '%s' to '%s'",
		source.Name, target.Name,
	)
}

func (s FromSchema) getSqlOrderBy(table string, query types.Query) []outputs.SqlOrderBy {
	var orderBy []outputs.SqlOrderBy
	for _, order := range query.OrderBy {
		orderBy = append(orderBy, outputs.SqlOrderBy{
			Table:     table,
			Attr:      order.Column.Name,
			Ascending: order.OrderBy == types.ASC,
		})
	}
	return orderBy
}

func (s FromSchema) getJoins(query types.Query, parentTable *types.Table, out *[]outputs.SqlJoinLine) error {
	for _, with := range query.With {
		line, err := s.getJoinLine(parentTable, with.Table)
		if err != nil {
			return err
		}

		line.Opts = outputs.SqlSelectOpts{
			Limit:   with.Query.Limit,
			Offset:  with.Query.Offset,
			Where:   with.Query.Where,
			OrderBy: s.getSqlOrderBy(with.Table.Name, with.Query),
		}

		*out = append(*out, line)

		err = s.getJoins(with.Query, with.Table, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s FromSchema) GetSelect(method types.Method) (outputs.SqlSelect, error) {
	var fields []outputs.SqlSelectField
	s.getSelectFields(method.Query, method.Table, &fields)

	var joins []outputs.SqlJoinLine
	err := s.getJoins(method.Query, method.Table, &joins)
	if err != nil {
		return outputs.SqlSelect{}, err
	}

	return outputs.SqlSelect{
		Table:  method.Table.Name,
		Select: fields,
		Joins:  joins,
		Opts: outputs.SqlSelectOpts{
			Limit:   method.Query.Limit,
			Offset:  method.Query.Offset,
			Where:   method.Query.Where,
			OrderBy: s.getSqlOrderBy(method.Table.Name, method.Query),
		},
	}, nil
}
