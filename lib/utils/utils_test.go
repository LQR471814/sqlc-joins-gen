package utils

import "testing"

func TestToGoIdentifier(t *testing.T) {
	testCases := []struct {
		test     string
		expected string
	}{
		{
			test:     "kebab-case-thing",
			expected: "Kebabcasething",
		},
		{
			test:     "snake_case",
			expected: "Snakecase",
		},
		{
			test:     "camelCase",
			expected: "Camelcase",
		},
		{
			test:     "PascalCase",
			expected: "Pascalcase",
		},
		{
			test:     "22-numbers-in-front",
			expected: "T22numbersinfront",
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
			return ToUpperGoIdentifier(c.test)
		}()
		if result != c.expected && c.expected != "<PANIC>" {
			t.Fatalf("expected '%s' but got '%s'", c.expected, result)
		}
	}
}
