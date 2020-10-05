package helpers

import (
	"strings"
	"unicode"
)

func KebabCase(s string) string {
	return strings.Trim(strings.Replace(SnakeCase(s), "_", "-", -1), "-")
}

func UpperSnakeCase(s string) string {
	return strings.ToUpper(SnakeCase(s))
}

func SnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return string(out)
}

func UpperCamelCase(s string) string {
	return camelCase(s, true)
}

func camelCase(s string, upper bool) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s))

	var prev rune
	for _, curr := range s {
		if !isDelimiter(curr) {
			if isDelimiter(prev) || (upper && prev == 0) {
				buffer = append(buffer, unicode.ToUpper(curr))
			} else {
				buffer = append(buffer, unicode.ToLower(curr))
			}
		}
		prev = curr
	}

	return string(buffer)
}

func isDelimiter(ch rune) bool {
	return ch == '-' || ch == '_' || unicode.IsSpace(ch)
}