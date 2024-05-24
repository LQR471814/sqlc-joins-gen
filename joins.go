package main

type JoinsList = []JoinQueryDef

type JoinQueryDef struct {
	Name   string        `json:"name"`
	Return ReturnMethods `json:"return"`
	Table  string        `json:"table"`
	Query  TableQuery    `json:"query"`
}

type TableQuery struct {
	Columns map[string]bool       `json:"columns,omitempty"`
	With    map[string]TableQuery `json:"with,omitempty"`
	Where   string                `json:"where,omitempty"`
	OrderBy map[string]string     `json:"orderBy,omitempty"`
	Limit   int                   `json:"limit,omitempty"`
	Offset  int                   `json:"offset,omitempty"`
}

type ReturnMethods = string

const (
	FIRST ReturnMethods = "first"
	MANY                = "many"
)
