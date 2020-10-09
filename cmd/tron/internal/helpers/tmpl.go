package helpers

import (
	"bytes"
	"text/template"

	"github.com/lissteron/simplerr"
)

func ExecTemplate(payload string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{"pkg": pkg(data)}).Parse(payload)
	if err != nil {
		return "", simplerr.Wrapf(err, "failed to parse template for %T", data)
	}

	buf := &bytes.Buffer{}

	if err := tmpl.Execute(buf, data); err != nil {
		return "", simplerr.Wrapf(err, "failed to execute template for %T", data)
	}

	return buf.String(), nil
}

func pkg(data interface{}) func(pkg string) string {
	if p, ok := data.(interface{ Pkg(string) string }); ok {
		return p.Pkg
	}

	return func(string) string { return "" }
}
