package generate

import (
	"go/format"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"golang.org/x/tools/imports"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Mainfile(p *project.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.CmdDir, p.Name, models.MainFile)

	printer.VerbosePrintln(color.FgMagenta, "Generate main.go")

	if !helpers.ConfirmOverwrite(path) {
		return nil
	}

	data := templates.NewMainData(&models.Data{Protos: p.Protos})

	data.AddImport("log")
	data.AddImport("github.com/loghole/tron")
	data.AddImport(strings.Join([]string{p.Module, "config"}, "/"))

	for _, proto := range p.Protos {
		data.AddImport(proto.Service.Import, proto.Service.Alias)
	}

	mainScript, err := helpers.ExecTemplate(templates.MainTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formattedBytes, err := format.Source([]byte(mainScript))
	if err != nil {
		return simplerr.Wrap(err, "failed to format process")
	}

	formattedBytes, err = imports.Process("", formattedBytes, nil)
	if err != nil {
		return simplerr.Wrap(err, "failed to imports process")
	}

	return helpers.WriteToFile(path, formattedBytes)
}
