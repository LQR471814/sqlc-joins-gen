package main

type CompositeQuery struct {
	Name   string               `json:"name"`
	Return CompositeQueryReturn `json:"return"`
	Table  string               `json:"table"`
	Query  CompositeQueryClause `json:"query"`
}

type CompositeQueryClause struct {
	Columns map[string]bool                 `json:"columns,omitempty"`
	With    map[string]CompositeQueryClause `json:"with,omitempty"`
	Where   string                          `json:"where,omitempty"`
	OrderBy map[string]string               `json:"orderBy,omitempty"`
	Limit   int                             `json:"limit,omitempty"`
	Offset  int                             `json:"offset,omitempty"`
}

type CompositeQueryReturn = string

const (
	FIRST CompositeQueryReturn = "first"
	MANY                       = "many"
)
