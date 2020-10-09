package generate

import (
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Makefile(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Makefile")

	data := templates.NewTronMKData(
		strings.Join([]string{models.CmdDir, p.Name, models.MainFile}, "/"),
		models.DockerfileFilepath,
		p.Module,
	)

	tronMK, err := helpers.ExecTemplate(templates.TronMK, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteToFile(models.TronMKFilepath, []byte(tronMK)); err != nil {
		return err
	}

	path := filepath.Join(p.AbsPath, models.MakefileFilepath)

	if err := helpers.WriteWithConfirm(path, []byte(templates.Makefile)); err != nil {
		return err
	}

	return nil
}
