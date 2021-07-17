package generate

import (
	"fmt"
	"go/format"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Config(project *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate config")

	config, err := helpers.ExecTemplate(templates.ConfigConstTemplate, project)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	formatted, err := format.Source([]byte(config))
	if err != nil {
		return fmt.Errorf("format sourse: %w", err)
	}

	path := filepath.Join(project.AbsPath, models.ConfigConstFilepath)

	if err := helpers.WriteToFile(path, formatted); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func ConfigHelper(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate config helper")

	path := filepath.Join(p.AbsPath, models.ConfigFilepath)

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(templates.ConfigTemplate)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
