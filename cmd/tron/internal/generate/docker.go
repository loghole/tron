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

func Dockerfile(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate Dockerfile")

	path := filepath.Join(p.AbsPath, models.DockerfileFilepath)

	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerfileTemplate, p)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteWithConfirm(path, []byte(dockerfile)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}
