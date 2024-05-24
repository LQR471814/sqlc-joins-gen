package main

import (
	"fmt"
	"strings"
)

func ToGoIdentifier(name string) string {
	result := ""

	for i, c := range name {
		if i == 0 {
			if c >= '0' && c <= '9' {
				result += "T"
			}
			result += strings.ToUpper(string(c))
			continue
		}

		if c == '-' || c == '_' {
			continue
		}

		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += strings.ToLower(string(c))
			continue
		}

		panic(fmt.Sprintf(
			"got invalid character in name to be converted to an identifier '%c'",
			c,
		))
	}

	return result
}
