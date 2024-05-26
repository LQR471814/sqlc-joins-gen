package querycfg

type Method struct {
	Name   string `json:"name"`
	Return Return `json:"return"`
	Table  string `json:"table"`
	Query  Query  `json:"query"`
}

type Query struct {
	Columns map[string]bool    `json:"columns,omitempty"`
	With    map[string]Query   `json:"with,omitempty"`
	Where   string             `json:"where,omitempty"`
	OrderBy map[string]OrderBy `json:"orderBy,omitempty"`
	Limit   int                `json:"limit,omitempty"`
	Offset  int                `json:"offset,omitempty"`
}

type Return = string

const (
	FIRST Return = "first"
	MANY         = "many"
)

type OrderBy = string

const (
	ASC OrderBy = "asc"
	DSC OrderBy = "dsc"
)
