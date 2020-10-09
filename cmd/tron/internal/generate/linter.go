package generate

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Linter(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate .golangci.yaml")

	path := filepath.Join(p.AbsPath, models.GolangciLintFilepath)

	if !helpers.ConfirmOverwrite(path) {
		return nil
	}

	if err := helpers.WriteToFile(path, []byte(templates.GolangCILintTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}
