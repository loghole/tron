package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func GoMod(p *models.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.GoModFile)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printer.VerbosePrintln(color.FgMagenta, "Initialise go mod")

		if err := helpers.Exec(p.AbsPath, "go", "mod", "init", p.Module); err != nil {
			return fmt.Errorf("exec cmd: %w", err)
		}

		printer.VerbosePrintln(color.FgBlue, "\tSuccess")

		return nil
	}

	return nil
}
