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

func Buf(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate buf.yaml and buf.gen.yaml")

	path := filepath.Join(p.AbsPath, models.BufFilepath)

	if helpers.ConfirmOverwrite(path) {
		if err := helpers.WriteToFile(path, []byte(templates.Buf)); err != nil {
			return fmt.Errorf("write file '%s': %w", path, err)
		}
	} else {
		printer.Println(color.FgBlue, "\tSkipped")
	}

	path = filepath.Join(p.AbsPath, models.BufGenFilepath)

	if helpers.ConfirmOverwrite(path) {
		if err := helpers.WriteToFile(path, []byte(templates.BufGen)); err != nil {
			return fmt.Errorf("write file '%s': %w", path, err)
		}
	} else {
		printer.Println(color.FgBlue, "\tSkipped")
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
