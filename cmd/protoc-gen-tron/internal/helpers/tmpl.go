package helpers

import (
	"bytes"
	"fmt"
	"text/template"
)

func ExecTemplate(payload string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{"pkg": pkg(data)}).Parse(payload)
	if err != nil {
		return "", fmt.Errorf("failed to parse template for %T: %w", data, err)
	}

	buf := &bytes.Buffer{}

	if err := tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template for %T: %w", data, err)
	}

	return buf.String(), nil
}

func pkg(data interface{}) func(pkg string) string {
	if p, ok := data.(interface{ Pkg(string) string }); ok {
		return p.Pkg
	}

	return func(string) string { return "" }
}
