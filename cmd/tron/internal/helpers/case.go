package helpers

import (
	"regexp"
	"strings"
	"unicode"
)

// nolint:gochecknoglobals //regexp
var (
	goNameRexp      = regexp.MustCompile(`[^a-zA-Z0-9\s_-]+`)
	firstDigitsRexp = regexp.MustCompile(`^\d+`)
)

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

func CamelCase(s string) string {
	return camelCase(s, false)
}

func GoName(s string) string {
	name := goNameRexp.ReplaceAllString(strings.ReplaceAll(s, ".", "_"), "")
	name = firstDigitsRexp.ReplaceAllString(name, "")

	return name
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
