package main

import (
	"flag"
	"log/slog"
	"os"
	"path"
	"sqlc-joins-gen/lib/inputs"
	"sqlc-joins-gen/lib/outputs"
)

func main() {
	sqlcFile := flag.String("config", "", "path to the sqlc.yaml config file")
	flag.Parse()

	isDir := true
	pathOrDir, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get current working dir", "err", err)
		return
	}
	if *sqlcFile != "" {
		pathOrDir = *sqlcFile
		isDir = false
	}

	tasks, err := inputs.LoadSqlcConfig(pathOrDir, isDir)
	if err != nil {
		slog.Error("failed to read sqlc config", "err", err)
		return
	}

	sqlgen := outputs.SqliteGenerator{}
	for _, task := range tasks {
		plgen := outputs.GolangGenerator{
			PackageName: task.Gen.Go.Package,
			PackagePath: path.Join(task.CfgDir, task.Gen.Go.Out),
		}
		err = generate(task, sqlgen, plgen)
		if err != nil {
			slog.Error("failed to execute sqlc task", "err", err)
		}
	}
}
