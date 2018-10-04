package linter

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func isCamelCase(s string) bool {
	if len(s) != 0 {
		c := s[0]
		return 'A' <= c && c <= 'Z' &&
			strings.IndexByte(s, '_') == -1 &&
			!isUpperUnderscore(s)
	}
	return false

}

func isCamelCase_X(s string) bool {
	first, _ := utf8.DecodeRuneInString(s)
	if unicode.IsLower(first) || s == strings.ToUpper(s) || strings.Contains(s, "_") {
		return false
	}
	return true
}

func isLowerUnderscore(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			return false
		}
	}
	return true
}

func isUpperUnderscore(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			return false
		}
	}
	return true
}
