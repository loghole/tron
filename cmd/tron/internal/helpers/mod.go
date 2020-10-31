package helpers

import (
	"strings"
)

func ModuleName(module string) string {
	parts := strings.Split(module, "/")

	return parts[len(parts)-1]
}
