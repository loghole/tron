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

func Git(p *models.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.GitDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printer.VerbosePrintln(color.FgMagenta, "Initialise git")

		if err := helpers.Exec(p.AbsPath, "git", "init"); err != nil {
			return fmt.Errorf("exec cmd: %w", err)
		}

		printer.VerbosePrintln(color.FgBlue, "\tSuccess")

		return nil
	}

	return nil
}
