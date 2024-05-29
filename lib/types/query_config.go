package types

import "fmt"

type MethodCfg struct {
	Name   string   `json:"name"`
	Table  string   `json:"table"`
	Return Return   `json:"return"`
	Query  QueryCfg `json:"query"`
}

func (cfg MethodCfg) ToMethod(schema Schema) (Method, error) {
	table := schema.FindTable(cfg.Table)
	if table == nil {
		return Method{}, fmt.Errorf("could not find table '%s'", cfg.Table)
	}
	query, err := cfg.Query.ToQuery(schema, table)
	if err != nil {
		return Method{}, err
	}
	return Method{
		Name:   cfg.Name,
		Table:  table,
		Return: cfg.Return,
		Query:  query,
	}, nil
}

type QueryCfg struct {
	Columns map[string]bool     `json:"columns,omitempty"`
	With    map[string]QueryCfg `json:"with,omitempty"`
	OrderBy map[string]OrderBy  `json:"orderBy,omitempty"`
	Where   string              `json:"where,omitempty"`
	Limit   int                 `json:"limit,omitempty"`
	Offset  int                 `json:"offset,omitempty"`
}

func (cfg QueryCfg) ToQuery(schema Schema, parentTable *Table) (Query, error) {
	query := Query{
		Where:  cfg.Where,
		Limit:  cfg.Limit,
		Offset: cfg.Offset,
	}

	for col, enabled := range cfg.Columns {
		column := parentTable.FindColumn(col)
		if column == nil {
			return Query{}, fmt.Errorf(
				"could not find column '%s' in table '%s'",
				col, parentTable.Name,
			)
		}
		query.Columns = append(query.Columns, QueryColumn{
			Column:  column,
			Enabled: enabled,
		})
	}

	for col, orderBy := range cfg.OrderBy {
		column := parentTable.FindColumn(col)
		if column == nil {
			return Query{}, fmt.Errorf(
				"could not find column '%s' in table '%s'",
				col, parentTable.Name,
			)
		}
		query.OrderBy = append(query.OrderBy, QueryOrderBy{
			Column:  column,
			OrderBy: orderBy,
		})
	}

	for childTableName, childQueryCfg := range cfg.With {
		childTable := schema.FindTable(childTableName)
		if childTable == nil {
			return Query{}, fmt.Errorf("could not find table '%s'", childTableName)
		}
		childQuery, err := childQueryCfg.ToQuery(schema, childTable)
		if err != nil {
			return Query{}, err
		}
		query.With = append(query.With, QueryWith{
			Table: childTable,
			Query: childQuery,
		})
	}

	return query, nil
}
