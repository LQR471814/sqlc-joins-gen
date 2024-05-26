package main

import (
	"flag"
	"log/slog"
	"os"
	"path"
	"sqlc-joins-gen/lib/gen"
	"sqlc-joins-gen/lib/querycfg"
	"sqlc-joins-gen/lib/sqlc"
	"sqlc-joins-gen/lib/sqlite"

	"github.com/titanous/json5"
)

func generate(task sqlc.CodegenTask) error {
	schema, err := sqlite.ParseSchema(task.Schema)
	if err != nil {
		return err
	}

	var methods []querycfg.Method
	err = json5.Unmarshal(task.Joins, &methods)
	if err != nil {
		return err
	}

	fromSchema := gen.GenManager{Schema: schema}
	return fromSchema.Generate(
		task,
		gen.SqliteGenerator{},
		gen.GolangGenerator{},
		methods,
	)
}

func main() {
	sqlcFile := flag.String("config", "", "path to the sqlc.yaml config file")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get current working dir", "err", err)
		return
	}
	if *sqlcFile != "" {
		dir = path.Dir(*sqlcFile)
	}

	tasks, err := sqlc.LoadConfig(dir)
	if err != nil {
		slog.Error("failed to read sqlc config", "err", err)
		return
	}

	for _, task := range tasks {
		err = generate(task)
		if err != nil {
			slog.Error("failed to execute sqlc task", "err", err)
		}
	}
}
