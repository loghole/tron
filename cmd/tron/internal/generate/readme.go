package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func ReadmeMD(p *models.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.ReadmeMDFilepath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		printer.Println(color.FgMagenta, "Generate README.md")

		template, err := helpers.ExecTemplate(templates.ReadmeMD, p)
		if err != nil {
			return fmt.Errorf("exec template: %w", err)
		}

		if err := helpers.WriteToFile(path, []byte(template)); err != nil {
			return fmt.Errorf("write file '%s': %w", path, err)
		}

		printer.Println(color.FgBlue, "\tSuccess")

		return nil
	}

	return nil
}
