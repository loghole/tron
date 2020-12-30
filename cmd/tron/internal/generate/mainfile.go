package generate

import (
	"fmt"
	"go/format"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/tools/imports"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Mainfile(project *models.Project, printer stdout.Printer) error {
	path := filepath.Join(project.AbsPath, models.CmdDir, project.Name, models.MainFile)

	printer.VerbosePrintln(color.FgMagenta, "Generate main.go")

	if !helpers.ConfirmOverwrite(path) {
		return nil
	}

	data := templates.NewMainData(project)

	for _, proto := range project.Protos {
		if !proto.WithImpl() {
			continue
		}

		dest := strings.Join([]string{
			project.Module,
			models.ProjectImportImplementation,
			strings.ReplaceAll(proto.Package, ".", "/"),
		}, "/")

		data.AddImport(dest, proto.GoPackage)
	}

	mainScript, err := helpers.ExecTemplate(templates.MainTemplate, data)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	formattedBytes, err := format.Source([]byte(mainScript))
	if err != nil {
		return fmt.Errorf("%s\n\nformat source: %w", mainScript, err)
	}

	formattedBytes, err = imports.Process("", formattedBytes, nil)
	if err != nil {
		return fmt.Errorf("%s\n\nimports process: %w", mainScript, err)
	}

	if err := helpers.WriteToFile(path, formattedBytes); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}
