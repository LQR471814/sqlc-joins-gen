package main

import (
	"fmt"
	"sqlc-joins-gen/lib/inputs"
	"sqlc-joins-gen/lib/outputs"
	"sqlc-joins-gen/lib/transform"
	"sqlc-joins-gen/lib/types"

	"github.com/titanous/json5"
)

func generate(
	task inputs.SqlcCodegenTask,
	sqlgen outputs.SqlGenerator,
	plgen outputs.PlGenerator,
) error {
	schema, err := inputs.ParseSqliteSchema(task.Schema)
	if err != nil {
		return err
	}

	var cfg []types.MethodCfg
	err = json5.Unmarshal(task.Joins, &cfg)
	if err != nil {
		return err
	}

	methods := make([]types.Method, len(cfg))
	for i, m := range cfg {
		methods[i], err = m.ToMethod(schema)
		if err != nil {
			return err
		}
	}

	fromSchema := transform.FromSchema{Schema: schema}

	script := outputs.PlScript{}
	for _, method := range methods {
		selectStmt, err := fromSchema.GetSelect(method)
		if err != nil {
			return err
		}

		out, err := fromSchema.GetMethodDef(method)
		if err != nil {
			return err
		}
		out.Sql = sqlgen.Select(selectStmt)

		for _, selectField := range selectStmt.Select {
			defIdx := -1
			for i, row := range out.RowDefs {
				if row.TableName == selectField.Table {
					defIdx = i
					break
				}
			}
			if defIdx < 0 {
				return fmt.Errorf(
					"could not find table '%s' in RowDefs for '%s'",
					selectField.Table, method.Name,
				)
			}

			fieldIdx := -1
			for i, field := range out.RowDefs[defIdx].Fields {
				if field.TableFieldName == selectField.Attr {
					fieldIdx = i
					break
				}
			}
			if fieldIdx < 0 {
				return fmt.Errorf(
					"could not find field '%s' in RowDef '%s' for '%s'",
					selectField.Table, out.RowDefs[defIdx].DefName, method.Name,
				)
			}

			out.ScanOrder = append(out.ScanOrder, outputs.PlScanEntry{
				RowDef: out.RowDefs[defIdx],
				Field:  out.RowDefs[defIdx].Fields[fieldIdx],
			})
		}

		out.RootDef = out.RowDefs[0]
		script.Methods = append(script.Methods, &out)
	}

	return plgen.Generate(script)
}
