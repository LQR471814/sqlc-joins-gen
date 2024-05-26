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

func (g GolangGenerator) typeStr(defs []PlRowDef, t PlType) string {
	if t.IsStruct {
		return defs[t.Struct].DefName
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

func (g GolangGenerator) writeStructDef(defs []PlRowDef, out *bytes.Buffer) {
	for _, row := range defs {
		out.WriteString(fmt.Sprintf(
			"// Table: %s\ntype %s struct {\n",
			row.TableName,
			upperGoIdentifier(row.DefName),
		))

		for _, def := range row.Fields {
			typeStr := g.typeStr(defs, def.Type)
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
	}
}

func (g GolangGenerator) scanRowCode(defs []PlRowDef, root PlRowDef, out *bytes.Buffer) {
	for _, field := range root.Fields {
		if field.Type.IsStruct {
			continue
		}
		out.WriteString(fmt.Sprintf(
			"&%s.%s,\n",
			root.DefName,
			upperGoIdentifier(field.Name),
		))
	}
	for _, field := range root.Fields {
		if !field.Type.IsStruct {
			continue
		}
		g.scanRowCode(defs, defs[field.Type.Struct], out)
	}
}

func (g GolangGenerator) queryFunc(defs []PlRowDef, root PlRowDef, out *bytes.Buffer) {
	out.WriteString(fmt.Sprintf(
		// TODO: add args, add single return
		"func (q *Queries) %s(ctx context.Context, args any) ([]%s, error) {\n",
		root.MethodName,
		root.MethodName,
	))

	out.WriteString(fmt.Sprintf(
		"rows, err := q.db.QueryContext(ctx, query%s, args)\n",
		root.MethodName,
	))
	out.WriteString("if err != nil { return nil, err }; defer rows.Close()\n\n")

	for _, row := range defs {
		out.WriteString(fmt.Sprintf("var %sMap map[string]%s\n", row.DefName, row.DefName))
	}

	out.WriteString("\nfor rows.Next() {\n")
	for _, row := range defs {
		out.WriteString(fmt.Sprintf("var %s %s\n", row.DefName, row.DefName))
	}
	out.WriteString("\nerr := rows.Scan(\n")
	g.scanRowCode(defs, root, out)
	out.WriteString(")\nif err != nil { return nil, err }\n\n")

	out.WriteString("}\n\n")
	out.WriteString(fmt.Sprintf("var items []%s\n", root.MethodName))
	out.WriteString(fmt.Sprintf("for _, i := range %sMap {\n", root.MethodName))
	out.WriteString("items = append(items, i)\n")
	out.WriteString("}\n")
	out.WriteString("if err := rows.Close(); err != nil { return nil, err }\n")
	out.WriteString("if err := rows.Err(); err != nil { return nil, err }\n")
	out.WriteString("return items, nil\n}")
}

func (g GolangGenerator) Script(cfg sqlc.CodegenTask, script PlScript) []PlScriptOutput {
	var sqlLocations []PlSqlLocation

	out := bytes.NewBufferString("// Code generated by sqlc-joins-gen. DO NOT EDIT.\n\n")
	out.WriteString(fmt.Sprintf("package %s\n\n", cfg.Gen.Go.Package))
	g.writeStructDef(script.RowDefs, out)

	for _, row := range script.RowDefs {
		if !row.MethodRoot {
			continue
		}
		out.WriteString(fmt.Sprintf("\n\nconst query%s = `", row.DefName))
		sqlLocations = append(sqlLocations, PlSqlLocation{
			MethodName: row.MethodName,
			Location:   out.Len(),
		})
		out.WriteString("`\n\n")
	}

	for _, row := range script.RowDefs {
		if !row.MethodRoot {
			continue
		}
		g.queryFunc(script.RowDefs, row, out)
	}

	return []PlScriptOutput{
		{
			Path:              path.Join(cfg.CfgDir, cfg.Gen.Go.Out, "query.joins.go"),
			Contents:          out.Bytes(),
			SqlEmbedLocations: sqlLocations,
		},
	}
}
