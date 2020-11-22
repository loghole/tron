package generate

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func TronMK(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Tron MK")

	data := templates.NewTronMKData(p)

	tronMK, err := helpers.ExecTemplate(templates.TronMK, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	path := filepath.Join(p.AbsPath, models.TronMKFilepath)

	if err := helpers.WriteToFile(path, []byte(tronMK)); err != nil {
		return err
	}

	return nil
}

func Makefile(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Makefile")

	path := filepath.Join(p.AbsPath, models.MakefileFilepath)

	if err := helpers.WriteWithConfirm(path, []byte(templates.Makefile)); err != nil {
		return err
	}

	return nil
}
