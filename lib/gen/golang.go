package gen

import (
	"bytes"
	"fmt"
	"path"
	"sqlc-joins-gen/lib/sqlc"
	"sqlc-joins-gen/lib/utils"
)

func lowerGoIdentifier(name string) string {
	result := ""

	for i, c := range name {
		if i == 0 && c >= '0' && c <= '9' {
			result += "T"
		}
		if c == '-' || c == '_' {
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += string(c)
			continue
		}
		panic(fmt.Sprintf(
			"got invalid character in name to be converted to an identifier '%c', '%s'",
			c,
			name,
		))
	}

	return result
}

func upperGoIdentifier(name string) string {
	return utils.Capitalize(lowerGoIdentifier(name))
}

type GolangGenerator struct{}

func (g GolangGenerator) typeStr(t PlType) string {
	if t.IsStruct {
		return t.Struct.Name
	}
	switch t.Primitive {
	case INT:
		if t.Nullable {
			return "sql.NullInt32"
		}
		return "int"
	case FLOAT:
		if t.Nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	case STRING:
		if t.Nullable {
			return "sql.NullString"
		}
		return "string"
	case BOOL:
		if t.Nullable {
			return "sql.NullBool"
		}
		return "bool"
	}
	panic(fmt.Sprintf("unknown primitive type '%d'", t.Primitive))
}

type GoStructField struct {
	Key  string
	Type string
}

type GoStructDef struct {
	StructName string
	TableName  string
	Fields     []GoStructField
}

func (g GolangGenerator) structDef(def PlRowDef) []GoStructDef {
	var defs []GoStructDef
	queue := []struct {
		structName string
		def        PlRowDef
	}{{def: def}}

	for len(queue) > 0 {
		row := queue[0]

		def := GoStructDef{
			StructName: upperGoIdentifier(row.structName),
			TableName:  row.def.Name,
		}
		for i, field := range row.def.Fields {
			typeStr := g.typeStr(field.Type)
			if field.Type.IsStruct {
				typeStr = def.StructName + fmt.Sprint(i)
				queue = append(queue, struct {
					structName string
					def        PlRowDef
				}{
					structName: typeStr,
					def:        field.Type.Struct,
				})
			}
			if field.Type.Array {
				typeStr = "[]" + typeStr
			}
			def.Fields = append(def.Fields, GoStructField{
				Key:  field.Name,
				Type: typeStr,
			})
		}
		defs = append(defs, def)

		if len(queue) > 0 {
			queue = queue[1:]
		}
	}

	return defs
}

func (g GolangGenerator) writeStructDef(out *bytes.Buffer, defs []GoStructDef) map[string]string {

	for len(queue) > 0 {
		row := queue[0]

		out.WriteString(fmt.Sprintf(
			"type %s struct {\n",
			upperGoIdentifier(row.Name),
		))

		for i, def := range row.Fields {
			typeStr := g.typeStr(def.Type)
			if def.Type.IsStruct {
				typeStr = row.Name + fmt.Sprint(i)
				tableToStruct[def.Type.Struct.Name] = typeStr
				queue = append(queue, PlRowDef{
					MethodName: row.MethodName,
					Name:       typeStr,
					Fields:     def.Type.Struct.Fields,
				})
			}
			if def.Type.Array {
				typeStr = "[]" + typeStr
			}

			out.WriteString(fmt.Sprintf(
				"%s %s\n",
				upperGoIdentifier(def.Name),
				typeStr,
			))
		}
		out.WriteString("}\n\n")

		if len(queue) > 0 {
			queue = queue[1:]
		}
	}

	return tableToStruct
}

func (g GolangGenerator) scanRowCode(rowDef PlRowDef, tableToStruct map[string]string, out *bytes.Buffer) {
	for _, field := range rowDef.Fields {
		if field.Type.IsStruct {
			continue
		}
		structName := tableToStruct[rowDef.Name]
		if structName == "" {
			structName = rowDef.Name
		}
		out.WriteString(fmt.Sprintf(
			"&%s.%s,\n",
			structName,
			upperGoIdentifier(field.Name),
		))
	}
	for _, field := range rowDef.Fields {
		if !field.Type.IsStruct {
			continue
		}
		g.scanRowCode(field.Type.Struct, tableToStruct, out)
	}
}

func (g GolangGenerator) queryFunc(rowDef PlRowDef, tableToStruct map[string]string, out *bytes.Buffer) {
	out.WriteString(fmt.Sprintf(
		// TODO: add args, add single return
		"func (q *Queries) %s(ctx context.Context, args any) ([]%s, error) {\n",
		rowDef.MethodName,
		rowDef.MethodName+"Row",
	))

	out.WriteString(fmt.Sprintf(
		"rows, err := q.db.QueryContext(ctx, query%sRow, args)\n",
		rowDef.MethodName,
	))
	out.WriteString("if err != nil { return nil, err }; defer rows.Close()\n\n")

	out.WriteString(fmt.Sprintf("var %sRowMap map[string]%sRow\n", rowDef.MethodName, rowDef.MethodName))
	for _, structName := range tableToStruct {
		out.WriteString(fmt.Sprintf("var %sMap map[string]%s\n", structName, structName))
	}

	out.WriteString("\nfor rows.Next() {\n")
	out.WriteString(fmt.Sprintf("var %sRow %sRow\n", rowDef.MethodName, rowDef.MethodName))
	for _, structName := range tableToStruct {
		out.WriteString(fmt.Sprintf("var %s %s\n", structName, structName))
	}
	out.WriteString("\nerr := rows.Scan(\n")
	g.scanRowCode(rowDef, tableToStruct, out)
	out.WriteString(")\nif err != nil { return nil, err }\n\n")

	out.WriteString("}\n\n")
	out.WriteString(fmt.Sprintf("var items []%s\n", rowDef.MethodName))
	out.WriteString(fmt.Sprintf("for _, i := range %sMap {\n", rowDef.MethodName))
	out.WriteString("items = append(items, i)\n")
	out.WriteString("}\n")
	out.WriteString("if err := rows.Close(); err != nil { return nil, err }\n")
	out.WriteString("if err := rows.Err(); err != nil { return nil, err }\n")
	out.WriteString("return items, nil\n}")
}

func (g GolangGenerator) Script(cfg sqlc.CodegenTask, script PlScript) []PlScriptOutput {
	var sqlLocations []PlSqlQueryLocation

	out := bytes.NewBufferString("// Code generated by sqlc-joins-gen. DO NOT EDIT.\n\n")
	out.WriteString(fmt.Sprintf("package %s\n\n", cfg.Gen.Go.Package))
	for _, d := range script.RowDefs {
		tableToStruct := g.structDef(out, d)

		out.WriteString(fmt.Sprintf("\n\nconst query%s = `", d.Name))
		sqlLocations = append(sqlLocations, PlSqlQueryLocation{
			MethodName: d.MethodName,
			Location:   out.Len(),
		})
		out.WriteString("`\n\n")

		g.queryFunc(d, tableToStruct, out)
	}

	return []PlScriptOutput{
		{
			Path:         path.Join(cfg.CfgDir, cfg.Gen.Go.Out, "query.joins.go"),
			Contents:     out.Bytes(),
			SqlLocations: sqlLocations,
		},
	}
}
