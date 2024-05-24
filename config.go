package main

import (
	"errors"
	"os"
	"path"
	"strings"
)

type SqlcConfig struct {
	Sql []SqlcTarget `yaml:"sql"`
}

type SqlcTarget struct {
	Engine  string `yaml:"engine"`
	Queries string `yaml:"queries"`
	Schema  string `yaml:"schema"`
	Gen     struct {
		Go struct {
			Package string `yaml:"package"`
			Out     string `yaml:"out"`
		} `yaml:"go"`
	} `yaml:"gen"`
}

type CodegenTask struct {
	Schema      []byte
	Joins       []byte
	PackageName string
	PackagePath string
}

func replaceExt(filename, ext string) string {
	lastDotIdx := strings.LastIndex(filename, ".")
	if lastDotIdx < 0 {
		return filename + "." + ext
	}
	return filename[:lastDotIdx] + "." + ext
}

func LoadConfig(cfgDir string, cfg SqlcConfig) ([]CodegenTask, error) {
	targets := []SqlcTarget{}
	for _, target := range cfg.Sql {
		if target.Engine == "sqlite" {
			targets = append(targets, target)
		}
	}
	if len(targets) == 0 {
		return nil, errors.New("no sqlc targets are of the sqlite engine")
	}

	tasks := []CodegenTask{}
	for _, target := range targets {
		schemaBuff, err := os.ReadFile(path.Join(cfgDir, target.Schema))
		if err != nil {
			return nil, err
		}

		joinsPath := path.Join(
			cfgDir,
			path.Dir(target.Queries),
			replaceExt(path.Base(target.Queries), "joins.json5"),
		)
		joinsBuff, err := os.ReadFile(joinsPath)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, CodegenTask{
			Schema:      schemaBuff,
			Joins:       joinsBuff,
			PackageName: target.Gen.Go.Package,
			PackagePath: target.Gen.Go.Out,
		})
	}

	return tasks, nil
}
