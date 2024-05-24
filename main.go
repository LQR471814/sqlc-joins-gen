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
	"gopkg.in/yaml.v3"
)

func generate(task sqlc.CodegenTask) error {
	schema, err := sqlite.ParseSchema(task.Schema)
	if err != nil {
		return err
	}

	var queries []querycfg.Query
	err = json5.Unmarshal(task.Joins, &queries)
	if err != nil {
		return err
	}

	fromSchema := gen.FromSchema{Schema: schema}
	for _, query := range queries {
		fromSchema.Generate(
			gen.SqliteGenerator{},
			gen.GolangGenerator{},
			query,
		)
	}

	return nil
}

func main() {
	sqlcFile := flag.String("config", "", "path to the sqlc.yaml config file")

	flag.Parse()

	if *sqlcFile == "" {
		slog.Error("-config must be specified")
		return
	}

	f, err := os.Open(*sqlcFile)
	if err != nil {
		slog.Error("failed to open sqlc config", "err", err)
		return
	}
	defer f.Close()

	sqlcCfg := sqlc.Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&sqlcCfg)
	if err != nil {
		slog.Error("failed to parse sqlc config", "err", err)
		return
	}

	tasks, err := sqlc.LoadConfig(path.Dir(*sqlcFile), sqlcCfg)
	if err != nil {
		slog.Error("failed to process sqlc config", "err", err)
		return
	}

	for _, task := range tasks {
		err = generate(task)
		if err != nil {
			slog.Error("failed to execute sqlc task", "err", err)
		}
	}
}
