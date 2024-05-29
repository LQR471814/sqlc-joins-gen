package types

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
	SourceColumn *Column
	TargetColumn *Column
}

type ForeignKey struct {
	TargetTable *Table
	On          []ForeignColumn
}

type Table struct {
	Name         string
	Columns      []*Column
	PrimaryKey   []*Column
	UniqueFields [][]*Column
	ForeignKeys  []ForeignKey
}

func (t Table) FindColumn(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (t Table) MustFindColumn(name string) *Column {
	res := t.FindColumn(name)
	if res == nil {
		panic(fmt.Sprintf("could not find column '%s'", name))
	}
	return res
}

type Schema struct {
	Tables []*Table
}

func (s Schema) FindTable(name string) *Table {
	for _, t := range s.Tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (s Schema) MustFindTable(name string) *Table {
	res := s.FindTable(name)
	if res == nil {
		panic(fmt.Sprintf("can't find table '%s'", name))
	}
	return res
}
