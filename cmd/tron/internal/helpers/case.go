package helpers

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	goNameRexp      = regexp.MustCompile(`[^a-zA-Z0-9\s_-]+`)
	firstDigitsRexp = regexp.MustCompile(`^\d+`)
)

func UpperCamelCase(s string) string {
	return camelCase(s, true)
}

func GoName(s string) string {
	name := goNameRexp.ReplaceAllString(strings.ReplaceAll(s, ".", "_"), "")
	name = firstDigitsRexp.ReplaceAllString(name, "")

	return name
}

func ProtoPkgName(s string) string {
	name := goNameRexp.ReplaceAllString(s, "")
	name = strings.ReplaceAll(name, "-", "_")

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
