package generate

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Linter(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate .golangci.yaml")

	path := filepath.Join(p.AbsPath, models.GolangciLintFilepath)

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	lint, err := helpers.ExecTemplate(templates.GolangCILintTemplate, p)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if err := helpers.WriteToFile(path, []byte(lint)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
