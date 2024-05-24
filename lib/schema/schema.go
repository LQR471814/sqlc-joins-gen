package schema

import (
	"fmt"
)

type ColumnType = string

const (
	TEXT ColumnType = "text"
	INT             = "integer"
	REAL            = "real"
)

type Column struct {
	Name     string
	Type     ColumnType
	Nullable bool
}

type ForeignColumn struct {
	SourceColumn int
	TargetColumn int
}

type ForeignKey struct {
	TargetTable int
	On          []ForeignColumn
}

type Table struct {
	Name        string
	Columns     []Column
	PrimaryKey  []int
	ForeignKeys []ForeignKey
}

func (t Table) FindColumnIdx(name string) int {
	for i, c := range t.Columns {
		if c.Name == name {
			return i
		}
	}
	return -1
}

func (t Table) MustFindColumnIdx(name string) int {
	idx := t.FindColumnIdx(name)
	if idx < 0 {
		panic(fmt.Sprintf("could not find column '%s'", name))
	}
	return idx
}

func (t Table) MustFindColumn(name string) Column {
	return t.Columns[t.MustFindColumnIdx(name)]
}

type Schema struct {
	Tables []Table
}

func (s Schema) FindTableIdx(name string) int {
	for i, t := range s.Tables {
		if t.Name == name {
			return i
		}
	}
	return -1
}

func (s Schema) MustFindTableIdx(name string) int {
	idx := s.FindTableIdx(name)
	if idx < 0 {
		panic(fmt.Sprintf("can't find table '%s'", name))
	}
	return idx
}

func (s Schema) MustFindTable(name string) Table {
	return s.Tables[s.MustFindTableIdx(name)]
}

