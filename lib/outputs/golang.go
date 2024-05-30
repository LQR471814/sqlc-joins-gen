package outputs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sqlc-joins-gen/lib/utils"
	"strings"
)

func goId(name string) string {
	result := ""

	for i, c := range name {
		if i == 0 && c >= '0' && c <= '9' {
			panic(fmt.Sprintf("identifier should not start with a number '%s'", name))
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

func upperGoId(id string) string {
	return utils.Capitalize(goId(id))
}

type GolangGenerator struct {
	PackageName string
	PackagePath string
}

func (g GolangGenerator) baseType(t PlType) string {
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

func (g GolangGenerator) typeStr(t PlType) string {
	str := g.baseType(t)
	if t.Array {
		return "[]" + str
	}
	return str
}

func (g GolangGenerator) writeStructDef(defs []*PlRowDef, out *bytes.Buffer) {
	for _, row := range defs {
		out.WriteString(fmt.Sprintf(
			"// Table: %s\ntype %s struct {\n",
			row.TableName, upperGoId(row.DefName),
		))

		for _, def := range row.Fields {
			var typeStr string
			if def.IsRowDef {
				typeStr = upperGoId(defs[def.RowDef].DefName)
				if def.Type.Array {
					typeStr = "[]" + typeStr
				}
			} else {
				typeStr = g.typeStr(def.Type)
			}
			out.WriteString(fmt.Sprintf(
				"%s %s\n",
				upperGoId(def.Name), typeStr,
			))
		}
		out.WriteString("}\n\n")
	}
}

func (g GolangGenerator) scanRowCode(defs []*PlRowDef, root *PlRowDef, out *bytes.Buffer) {
	for _, field := range root.Fields {
		if field.IsRowDef {
			continue
		}
		out.WriteString(fmt.Sprintf(
			"&%s.%s,\n",
			goId(root.DefName),
			upperGoId(field.Name),
		))
	}
	for _, field := range root.Fields {
		if !field.IsRowDef {
			continue
		}
		g.scanRowCode(defs, defs[field.RowDef], out)
	}
}

func (g GolangGenerator) queryFunc(
	args []PlQueryArg,
	method *PlMethodDef,
	out *bytes.Buffer,
) {
	defs := method.RowDefs
	upperRootDefName := utils.Capitalize(method.RootDef.DefName)

	var returnType string
	if method.FirstOnly {
		returnType = "*" + upperRootDefName
	} else {
		returnType = fmt.Sprintf("[]%s", upperRootDefName)
	}
	errReturnStmt := "if err != nil { return nil, err }"

	var argsDef string
	for i, arg := range args {
		if i > 0 {
			argsDef += ", "
		}
		argsDef += fmt.Sprintf("%s %s", goId(arg.Name), g.typeStr(arg.Type))
	}

	out.WriteString(fmt.Sprintf(
		"func (q *Queries) %s(ctx context.Context, %s) (%s, error) {\n",
		utils.Capitalize(method.MethodName), argsDef, returnType,
	))

	out.WriteString(fmt.Sprintf("queryStr := %sQuery\n", method.RootDef.DefName))

	for i, arg := range args {
		out.WriteString("\n")
		argName := goId(arg.Name)

		formatStmt := ""
		switch arg.Type.Primitive {
		case INT:
			if arg.Type.Nullable {
				formatStmt += fmt.Sprintf("%sStr := \"null\"\n", argName)
				formatStmt += fmt.Sprintf("if %s.Valid {\n", argName)
				formatStmt += fmt.Sprintf("%sStr = fmt.Sprint(__ARG__.Int32)\n", argName)
				formatStmt += "}\n"
				break
			}
			formatStmt = fmt.Sprintf("%sStr := fmt.Sprint(__ARG__)\n", argName)
		case FLOAT:
			if arg.Type.Nullable {
				formatStmt += fmt.Sprintf("%sStr := \"null\"\n", argName)
				formatStmt += "if __ARG__.Valid {\n"
				formatStmt += fmt.Sprintf("%sStr = fmt.Sprint(__ARG__.Float64)\n", argName)
				formatStmt += "}\n"
				break
			}
			formatStmt = fmt.Sprintf("%sStr := fmt.Sprint(__ARG__)\n", argName)
		case STRING:
			if arg.Type.Nullable {
				formatStmt += fmt.Sprintf("%sStr := \"null\"\n", argName)
				formatStmt += "if __ARG__.Valid {\n"
				formatStmt += fmt.Sprintf("%sStr = `\"` + __ARG__.String + `\"`\n", argName)
				formatStmt += "}\n"
				break
			}
			formatStmt = fmt.Sprintf("%sStr := `\"` + __ARG__ + `\"`\n", argName)
		case BOOL:
			if arg.Type.Nullable {
				formatStmt += fmt.Sprintf("%sStr := \"null\"\n", argName)
				formatStmt += "if __ARG__.Valid {\n"
				formatStmt += fmt.Sprintf("%sInt := 0\n", argName)
				formatStmt += fmt.Sprintf("if __ARG__.Bool {\n")
				formatStmt += fmt.Sprintf("%sInt = 1\n", argName)
				formatStmt += "}\n"
				formatStmt += fmt.Sprintf("%sStr = fmt.Sprint(%sInt)\n", argName, argName)
				formatStmt += "}\n"
				break
			}
			formatStmt += fmt.Sprintf("%sInt := 0\n", argName)
			formatStmt += "if __ARG__ {\n"
			formatStmt += fmt.Sprintf("%sInt = 1\n", argName)
			formatStmt += "}\n"
			formatStmt += fmt.Sprintf("%sStr := fmt.Sprint(%sInt)\n", argName, argName)
		}

		if arg.Type.Array {
			out.WriteString(fmt.Sprintf("%sJoined := \"\"\n", argName))
			out.WriteString(fmt.Sprintf("for i, e := range %s {\n", argName))
			out.WriteString(fmt.Sprintf("if i > 0 { %sJoined += \", \" }\n", argName))
			out.WriteString(strings.ReplaceAll(formatStmt, "__ARG__", "e"))
			out.WriteString(fmt.Sprintf("%sJoined += %sStr\n", argName, argName))
			out.WriteString("}\n")
			out.WriteString(fmt.Sprintf(
				"queryStr = strings.Replace(queryStr, \"$%d\", %sJoined, 1)\n",
				i, argName,
			))
			continue
		}
		out.WriteString(strings.ReplaceAll(formatStmt, "__ARG__", argName))
		out.WriteString(fmt.Sprintf(
			"queryStr = strings.Replace(queryStr, \"$%d\", %sStr, 1)\n",
			i, argName,
		))
	}

	out.WriteString("\nrows, err := q.db.QueryContext(ctx, queryStr)\n")
	out.WriteString(fmt.Sprintf("%s; defer rows.Close()\n\n", errReturnStmt))

	for _, row := range defs {
		row.DefName = goId(row.DefName)
		out.WriteString(fmt.Sprintf("%sMap := newQueryMap[%s]()\n", row.DefName, utils.Capitalize(row.DefName)))
	}

	out.WriteString("\nfor rows.Next() {\n")
	for _, row := range defs {
		out.WriteString(fmt.Sprintf("var %s %s\n", row.DefName, utils.Capitalize(row.DefName)))
	}
	out.WriteString("\nerr := rows.Scan(\n")
	g.scanRowCode(defs, method.RootDef, out)
	out.WriteString(fmt.Sprintf(")\n%s", errReturnStmt))

	for _, row := range defs {
		out.WriteString("\n\n")

		out.WriteString(fmt.Sprintf("%sPkey := fmt.Sprint(", row.DefName))
		for i, col := range row.PrimaryKey {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(fmt.Sprintf("%s.%s", row.DefName, utils.Capitalize(col.Name)))
		}
		out.WriteString(")\n")

		out.WriteString(fmt.Sprintf(
			"%sExisting, ok := %sMap.dict[%sPkey]\n",
			row.DefName, row.DefName, row.DefName,
		))
		out.WriteString("if !ok {\n")
		out.WriteString(fmt.Sprintf(
			"%sMap.list = append(%sMap.list, %s)\n",
			row.DefName, row.DefName, row.DefName,
		))
		out.WriteString(fmt.Sprintf(
			"%sMap.dict[%sPkey] = &%sMap.list[len(%sMap.list)-1]\n",
			row.DefName, row.DefName, row.DefName, row.DefName,
		))
		if row.Parent != nil && row.ParentField != nil {
			if row.ParentField.Type.Array {
				out.WriteString(fmt.Sprintf(
					"%sExisting.%s = append(%sExisting.%s, *%sExisting)\n",
					row.Parent.DefName, row.TableName, row.Parent.DefName, row.TableName,
					row.DefName,
				))
			} else {
				out.WriteString(fmt.Sprintf(
					"%sExisting.%s = *%sExisting\n",
					row.Parent.DefName, row.TableName, row.DefName,
				))
			}
		}

		out.WriteString("}")
	}

	out.WriteString("}\n\n")
	out.WriteString("err = rows.Close()\n")
	out.WriteString(fmt.Sprintf("%s\n", errReturnStmt))
	out.WriteString("err = rows.Err()\n")
	out.WriteString(fmt.Sprintf("%s\n", errReturnStmt))

	if !method.FirstOnly {
		out.WriteString(fmt.Sprintf("return %sMap.list, nil\n}", method.MethodName))
		return
	}
	out.WriteString(fmt.Sprintf("if len(%sMap.list) == 0 { return nil, nil }\n", method.MethodName))
	out.WriteString(fmt.Sprintf("return &%sMap.list[0], nil\n}", method.MethodName))
}

type GoCodegenOptions struct {
	PackageName string
	PackagePath string
}

func (g GolangGenerator) Generate(script PlScript) error {
	out := bytes.NewBufferString("// Code generated by sqlc-joins-gen. DO NOT EDIT.\n\n")
	out.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))

	out.WriteString("import \"github.com/petar/GoLLRB/llrb\"\n\n")
	out.WriteString(`type queryMap[T any] struct {
	dict map[string]*T
	list []T
}

func newQueryMap[T any]() queryMap[T] {
	return queryMap[T]{
		dict: make(map[string]*T),
	}
}`)
	out.WriteString("\n\n")

	for _, method := range script.Methods {
		g.writeStructDef(method.RowDefs, out)
		out.WriteString(fmt.Sprintf(
			"\n\nconst %sQuery = `%s`\n\n",
			method.RootDef.DefName,
			method.Sql,
		))
		g.queryFunc(method.Args, method, out)
	}

	err := os.MkdirAll(g.PackagePath, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}
	outPath := path.Join(g.PackagePath, "query.joins.go")
	err = os.WriteFile(outPath, out.Bytes(), 0777)
	if err != nil {
		return err
	}
	err = os.Chdir(path.Dir(outPath))
	if err != nil {
		return err
	}
	err = exec.Command("go", "fmt").Run()
	if err != nil {
		return err
	}
	withImports, err := exec.Command("goimports", outPath).Output()
	if err != nil {
		return fmt.Errorf("failed to resolve imports with `goimports`, is it installed?\n'%s'", err.Error())
	}
	return os.WriteFile(outPath, []byte(withImports), 0777)
}
