package transform

import (
	"fmt"
	"github.com/lqr471814/sqlc-joins-gen/lib/outputs"
	"github.com/lqr471814/sqlc-joins-gen/lib/types"
)

// convert an sql column type into a primitive type
func SqlColumnTypeToPlType(t types.ColumnType) outputs.PlPrimitive {
	switch t {
	case types.INT:
		return outputs.INT
	case types.TEXT:
		return outputs.STRING
	case types.REAL:
		return outputs.FLOAT
	}
	panic(fmt.Sprintf("unknown column type '%s'", t))
}

