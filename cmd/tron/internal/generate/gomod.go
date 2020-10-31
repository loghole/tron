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

func GoMod(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Initialise go mod")

	path := filepath.Join(p.AbsPath, models.GoModFile)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return helpers.Exec(p.AbsPath, "go", "mod", "init", p.Module)
		}

		return err
	}

	return nil
}
