package gen

import (
	"fmt"
	"sqlc-joins-gen/lib/querycfg"
	"sqlc-joins-gen/lib/schema"
)

type FromSchema struct {
	Schema schema.Schema
}

func (from FromSchema) getColumns(qclause querycfg.Clause, table schema.Table) []schema.Column {
	columns := []schema.Column{}
	for col, enabled := range qclause.Columns {
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

func (from FromSchema) getPlStructDef(qclause querycfg.Clause, table schema.Table) PlStructDef {
	def := PlStructDef{
		Name: table.Name + "Row",
	}

	columns := from.getColumns(qclause, table)
	for _, col := range columns {
		def.Fields.Fields = append(def.Fields.Fields, PlFieldDef{
			Name: col.Name,
			Type: PlType{
				Primitive: SqlColumnTypeToPlType(col.Type),
			},
			Nullable: col.Nullable,
		})
	}

	for childTableName, childClause := range qclause.With {
		childTable := from.Schema.Tables[from.Schema.MustFindTableIdx(childTableName)]
		def.Fields.Fields = append(def.Fields.Fields, PlFieldDef{
			Name: childTableName,
			Type: PlType{
				IsStruct: true,
				Struct:   from.getPlStructDef(childClause, childTable),
			},
		})
	}

	return def
}

func (from FromSchema) getSelectFields(qclause querycfg.Clause, table schema.Table, out *[]SqlSelectField) {
	columns := from.getColumns(qclause, table)

	for _, col := range columns {
		*out = append(*out, SqlSelectField{
			Table: table.Name,
			Attr:  col.Name,
			As:    fmt.Sprintf("%s_%s", table.Name, col.Name),
		})
	}

	for tableName, childQuery := range qclause.With {
		table := from.Schema.Tables[from.Schema.MustFindTableIdx(tableName)]
		from.getSelectFields(childQuery, table, out)
	}
}

func (from FromSchema) getJoinLines(sourceIdx, targetIdx int, out *[]SqlJoinLine) {
	source := from.Schema.Tables[sourceIdx]
	target := from.Schema.Tables[targetIdx]

	for _, fkey := range source.ForeignKeys {
		if fkey.TargetTable == targetIdx {
			line := SqlJoinLine{Table: target.Name}
			for _, col := range fkey.On {
				line.On = append(line.On, SqlJoinOn{
					SourceTable: source.Name,
					SourceAttr:  source.Columns[col.SourceColumn].Name,
					TargetTable: target.Name,
					TargetAttr:  target.Columns[col.TargetColumn].Name,
				})
			}
			*out = append(*out, line)
			return
		}
	}

	for _, fkey := range target.ForeignKeys {
		if fkey.TargetTable == sourceIdx {
			line := SqlJoinLine{Table: target.Name}
			for _, col := range fkey.On {
				line.On = append(line.On, SqlJoinOn{
					SourceTable: target.Name,
					SourceAttr:  target.Columns[col.TargetColumn].Name,
					TargetTable: source.Name,
					TargetAttr:  source.Columns[col.SourceColumn].Name,
				})
			}
			*out = append(*out, line)
			return
		}
	}
}

func (from FromSchema) getJoins(qclause querycfg.Clause, parentTableIdx int, out *[]SqlJoinLine) {
	for childTable, childQuery := range qclause.With {
		targetTableIdx := from.Schema.MustFindTableIdx(childTable)
		from.getJoinLines(parentTableIdx, targetTableIdx, out)
		from.getJoins(childQuery, targetTableIdx, out)
	}
}

func (from FromSchema) getSelect(query querycfg.Query) SqlSelect {
	fields := []SqlSelectField{}
	tableIdx := from.Schema.MustFindTableIdx(query.Table)
	table := from.Schema.Tables[tableIdx]
	from.getSelectFields(query.Clause, table, &fields)

	joins := []SqlJoinLine{}
	from.getJoins(query.Clause, tableIdx, &joins)

	return SqlSelect{
		Table:  query.Name,
		Select: SqlSelectClause{Fields: fields},
		Joins:  SqlJoinClause{Joins: joins},
	}
}

func (from FromSchema) Generate(
	sqlgen SqlGenerator,
	plgen PlGenerator,
	query querycfg.Query,
) {
	selectStmt := from.getSelect(query)
	structDef := from.getPlStructDef(
		query.Clause,
		from.Schema.Tables[from.Schema.MustFindTableIdx(query.Table)],
	)
	fmt.Println(sqlgen.Select(selectStmt))
	fmt.Println(plgen.StructDef(structDef))
}
