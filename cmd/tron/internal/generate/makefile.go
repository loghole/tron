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

func TronMK(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Tron MK")

	tronMK, err := helpers.ExecTemplate(templates.TronMK, project)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	path := filepath.Join(project.AbsPath, models.TronMKFilepath)

	if err := helpers.WriteToFile(path, []byte(tronMK)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func Makefile(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Makefile")

	path := filepath.Join(project.AbsPath, models.MakefileFilepath)

	if !helpers.ConfirmOverwrite(path) {
		printer.VerbosePrintln(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(templates.Makefile)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}
