package utils

import (
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Capitalize(text string) string {
	if len(text) == 0 {
		return ""
	}
	if len(text) == 1 {
		return strings.ToUpper(text[0:1])
	}
	return strings.ToUpper(text[0:1]) + text[1:]
}

func lessStr(a, b string) bool {
	return a < b
}
func lessInt(a, b int) bool {
	return a < b
}
func lessFloat(a, b float32) bool {
	return a < b
}

func DiffUnordered(expected, got any) string {
	return cmp.Diff(
		expected, got,
		cmpopts.SortSlices(lessStr),
		cmpopts.SortSlices(lessInt),
		cmpopts.SortSlices(lessFloat),
		cmpopts.SortMaps(lessStr),
		cmpopts.SortMaps(lessInt),
		cmpopts.SortMaps(lessFloat),
	)
}
