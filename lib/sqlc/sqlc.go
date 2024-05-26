package sqlc

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sql []Target `yaml:"sql"`
}

type GenCfg struct {
	Go struct {
		Package string `yaml:"package"`
		Out     string `yaml:"out"`
	} `yaml:"go"`
}

type Target struct {
	Engine  string `yaml:"engine"`
	Queries string `yaml:"queries"`
	Schema  string `yaml:"schema"`
	Gen     GenCfg `yaml:"gen"`
}

type CodegenTask struct {
	CfgDir string
	Schema []byte
	Joins  []byte
	Gen    GenCfg
}

func replaceExt(filename, ext string) string {
	lastDotIdx := strings.LastIndex(filename, ".")
	if lastDotIdx < 0 {
		return filename + "." + ext
	}
	return filename[:lastDotIdx] + "." + ext
}

func readSqlcConfig(pathOrDir string, isDir bool) (*os.File, error) {
	if !isDir {
		return os.Open(pathOrDir)
	}
	f, err := os.Open(path.Join(pathOrDir, "sqlc.yaml"))
	if err == nil {
		return f, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	f, err = os.Open(path.Join(pathOrDir, "sqlc.yml"))
	if err == nil {
		return f, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	return nil, fmt.Errorf(
		"could not find sqlc.yaml or sqlc.yml in '%s'",
		pathOrDir,
	)
}

func LoadConfig(pathOrDir string, isDir bool) ([]CodegenTask, error) {
	f, err := readSqlcConfig(pathOrDir, isDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dir := pathOrDir
	if !isDir {
		dir = path.Dir(pathOrDir)
	}

	cfg := Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	var targets []Target
	for _, target := range cfg.Sql {
		if target.Engine == "sqlite" {
			targets = append(targets, target)
		}
	}
	if len(targets) == 0 {
		return nil, errors.New("no sqlc targets are of the sqlite engine")
	}

	var tasks []CodegenTask
	for _, target := range targets {
		schemaBuff, err := os.ReadFile(path.Join(dir, target.Schema))
		if err != nil {
			return nil, err
		}

		joinsPath := path.Join(
			dir,
			path.Dir(target.Queries),
			replaceExt(path.Base(target.Queries), "joins.json5"),
		)
		joinsBuff, err := os.ReadFile(joinsPath)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, CodegenTask{
			CfgDir: dir,
			Schema: schemaBuff,
			Joins:  joinsBuff,
			Gen:    target.Gen,
		})
	}

	return tasks, nil
}
