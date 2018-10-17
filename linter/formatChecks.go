package linter

func isLower(c byte) bool { return 'a' <= c && c <= 'z' }
func isUpper(c byte) bool { return 'A' <= c && c <= 'Z' }
func isDigit(c byte) bool { return '0' <= c && c <= '9' }

func isCamelCase(s string) bool {
	const minAllUppercaseLen = 4
	if len(s) == 0 || isLower(s[0]) {
		return false
	}
	i := 1
	for ; i < len(s); i++ {
		c := s[i]
		if isLower(c) {
			break
		}
		if !isUpper(c) && !isDigit(c) {
			return false
		}
	}
	if n := len(s); i == n {
		return n <= 4
	}
	for ; i < len(s); i++ {
		c := s[i]
		if !isLower(c) && !isUpper(c) && !isDigit(c) {
			return false
		}
	}
	return true
}

func isLowerUnderscore(s string) bool {
	if len(s) == 0 || !isLower(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !isLower(c) && !isDigit(c) && c != '_' {
			return false
		}
	}
	return isLower(s[len(s)-1]) // don't allow trailing underscore
}

func isUpperUnderscore(s string) bool {
	if len(s) == 0 || !isUpper(s[0]) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !isUpper(c) && !isDigit(c) && c != '_' {
			return false
		}
	}
	return isUpper(s[len(s)-1]) // don't allow trailing underscore
}
