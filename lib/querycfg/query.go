package querycfg

type Query struct {
	Name   string `json:"name"`
	Return Return `json:"return"`
	Table  string `json:"table"`
	Clause Clause `json:"query"`
}

type Clause struct {
	Columns map[string]bool   `json:"columns,omitempty"`
	With    map[string]Clause `json:"with,omitempty"`
	Where   string            `json:"where,omitempty"`
	OrderBy map[string]string `json:"orderBy,omitempty"`
	Limit   int               `json:"limit,omitempty"`
	Offset  int               `json:"offset,omitempty"`
}

type Return = string

const (
	FIRST Return = "first"
	MANY         = "many"
)
