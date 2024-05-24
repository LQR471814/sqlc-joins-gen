package main

import (
	"flag"
	"log/slog"
	"os"
	"path"

	"github.com/titanous/json5"
	"gopkg.in/yaml.v3"
)

func generate(task CodegenTask) error {
	schema, err := NewSchema(task.Schema)
	if err != nil {
		return err
	}

	var joins JoinsList
	err = json5.Unmarshal(task.Joins, &joins)
	if err != nil {
		return err
	}

	generator := JoinsGenerator{
		Schema: schema,
		Joins:  joins,
	}
	generator.Generate()

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

	sqlcCfg := SqlcConfig{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&sqlcCfg)
	if err != nil {
		slog.Error("failed to parse sqlc config", "err", err)
		return
	}

	tasks, err := LoadConfig(path.Dir(*sqlcFile), sqlcCfg)
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
