package generate

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func ReadmeMD(p *project.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.ReadmeMDFilepath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printer.VerbosePrintln(color.FgMagenta, "Generate README.md")

		template, err := helpers.ExecTemplate(templates.ReadmeMD, p)
		if err != nil {
			return simplerr.Wrap(err, "failed to exec template")
		}

		if err := helpers.WriteWithConfirm(path, []byte(template)); err != nil {
			return simplerr.Wrap(err, "failed to write file")
		}

		return nil
	}

	return ErrAlreadyExists
}
