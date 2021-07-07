package generator

import (
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
)

const appPath = "internal/app"

func (gen *Generator) implDir(pkg string) string {
	p := strings.ReplaceAll(
		strings.TrimPrefix(pkg, gen.protoPkgPrefix),
		".",
		string(filepath.Separator),
	)

	p = filepath.Join(
		appPath,
		p,
	)

	return p
}

func (gen *Generator) implImport(pkg string) protogen.GoImportPath {
	p := strings.ReplaceAll(
		strings.TrimPrefix(pkg, gen.protoPkgPrefix),
		".",
		"/",
	)

	p = path.Join(
		gen.module,
		appPath,
		p,
	)

	p = strings.TrimPrefix(p, "/")

	return protogen.GoImportPath(p)
}

func fileName(s string) string {
	f := filepath.Base(s)
	ext := filepath.Ext(f)

	return f[:len(f)-len(ext)]
}

func snakeCase(s string) string {
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
