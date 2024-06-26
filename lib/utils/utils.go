package utils

import (
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/exp/constraints"
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

func less[T constraints.Ordered](a, b T) bool {
	return a < b
}

func primitive() []cmp.Option {
	return []cmp.Option{
		cmpopts.SortSlices(less[string]),
		cmpopts.SortSlices(less[int]),
		cmpopts.SortSlices(less[int32]),
		cmpopts.SortSlices(less[int64]),
		cmpopts.SortSlices(less[uint]),
		cmpopts.SortSlices(less[uint32]),
		cmpopts.SortSlices(less[uint64]),
		cmpopts.SortSlices(less[float32]),
		cmpopts.SortSlices(less[float64]),
	}
}

var sorts = primitive()

func AddCustomSort[T any](newSorts ...func(val T) string) {
	for _, fn := range newSorts {
		res := func(a, b T) bool {
			return fn(a) < fn(b)
		}
		sorts = append(sorts, cmpopts.SortSlices(res))
	}
}

func DiffUnordered(expected, got any, opts ...cmp.Option) string {
	return cmp.Diff(expected, got, append(sorts, opts...)...)
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	cursor := 0
	for _, match := range re.FindAllStringSubmatchIndex(str, -1) {
		groups := []string{}
		for i := 0; i < len(match); i += 2 {
			if match[i] < 0 {
				groups = append(groups, "")
				continue
			}
			groups = append(groups, str[match[i]:match[i+1]])
		}
		matchStart := match[0]
		matchEnd := match[1]
		result += str[cursor:matchStart] + repl(groups)
		cursor = matchEnd
	}
	return result + str[cursor:]
}
