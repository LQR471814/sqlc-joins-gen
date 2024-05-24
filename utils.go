package main

import (
	"fmt"
	"strings"
)

func Capitalize(text string) string {
	if len(text) == 0 {
		return ""
	}
	if len(text) == 1 {
		return strings.ToUpper(text[0:1])
	}
	return strings.ToUpper(text[0:1]) + strings.ToLower(text[1:])
}

func ToLowerGoIdentifier(name string) string {
	result := ""

	for i, c := range name {
		if i == 0 && c >= '0' && c <= '9' {
			result += "T"
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

func ToUpperGoIdentifier(name string) string {
	return Capitalize(ToLowerGoIdentifier(name))
}

func SqliteToGoType(t ColumnType) string {
	switch t {
	case TEXT:
		return "string"
	case INT:
		return "int"
	case REAL:
		return "float64"
	}
	panic(fmt.Sprintf("got invalid sqlite type '%s'", t))
}
