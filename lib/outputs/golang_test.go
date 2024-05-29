package outputs

import "testing"

func TestToGoIdentifier(t *testing.T) {
	testCases := []struct {
		test     string
		expected string
	}{
		{
			test:     "kebab-case-thing",
			expected: "kebabcasething",
		},
		{
			test:     "snake_case",
			expected: "snakecase",
		},
		{
			test:     "camelCase",
			expected: "camelCase",
		},
		{
			test:     "PascalCase",
			expected: "PascalCase",
		},
		{
			test:     "22-numbers-in-front",
			expected: "<PANIC>",
		},
		{
			test:     "invalid thing with spaces",
			expected: "<PANIC>",
		},
	}

	for _, c := range testCases {
		result := func() string {
			defer func() {
				err := recover()
				if err == nil {
					return
				}
				if c.expected == "<PANIC>" {
					return
				}
				panic(err)
			}()
			return goId(c.test)
		}()
		if result != c.expected && c.expected != "<PANIC>" {
			t.Fatalf("expected '%s' but got '%s'", c.expected, result)
		}
	}
}

