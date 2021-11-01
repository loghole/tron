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

func MainFile(project *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate main.go")

	data, err := helpers.ExecTemplate(templates.MainFile, project)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	path := filepath.Join(project.AbsPath, models.CmdDir, project.Name, models.MainFile)

	if err := helpers.WriteToFile(path, []byte(data)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
