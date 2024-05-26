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

func less(a, b string) bool {
	return a < b
}

func DiffUnordered(expected, got any) string {
	return cmp.Diff(
		expected, got,
		cmpopts.SortSlices(less),
		cmpopts.SortMaps(less),
	)
}
