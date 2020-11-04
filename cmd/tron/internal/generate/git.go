package generate

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func Git(p *project.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.GitDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printer.VerbosePrintln(color.FgMagenta, "Initialise git")

		return helpers.Exec(p.AbsPath, "git", "init")
	}

	return ErrAlreadyExists
}
