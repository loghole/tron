package generate

import (
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Gitignore(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgCyan, "Generate .gitignore")

	path := filepath.Join(p.AbsPath, models.GitignoreFilepath)

	if !helpers.ConfirmOverwrite(path) {
		return nil
	}

	return helpers.WriteToFile(path, []byte(templates.GitignoreTemplate))
}
