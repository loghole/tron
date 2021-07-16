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

func Gitignore(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate .gitignore")

	path := filepath.Join(p.AbsPath, models.GitignoreFilepath)

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(templates.GitignoreTemplate)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
