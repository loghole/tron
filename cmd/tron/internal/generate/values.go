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
	"github.com/loghole/tron/internal/app"
)

func Values(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate values")

	data := templates.ValuesData{
		List: []templates.Env{
			{Key: strings.ToUpper(app.LoggerLevelEnv), Val: `info`},
			{Key: strings.ToUpper(app.LoggerCollectorAddrEnv), Val: `""`},
			{Key: strings.ToUpper(app.LoggerDisableStdoutEnv), Val: `false`},
			{Key: strings.ToUpper(app.JaegerAddrEnv), Val: `""`},
			{Key: strings.ToUpper(app.HTTPPortEnv), Val: `8080`},
			{Key: strings.ToUpper(app.GRPCPortEnv), Val: `8081`},
			{Key: strings.ToUpper(app.AdminPortEnv), Val: `8082`},
		},
	}

	values, err := helpers.ExecTemplate(templates.ValuesTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	writers := []struct {
		path string
		data string
	}{
		{models.ValuesBaseFilepath, values},
		{models.ValuesDevFilepath, templates.ValuesDevTemplate},
		{models.ValuesLocalFilepath, templates.ValuesLocalTemplate},
		{models.ValuesStgFilepath, templates.ValuesStgTemplate},
		{models.ValuesProdFilepath, templates.ValuesProdTemplate},
	}

	for _, wr := range writers {
		if err := helpers.WriteWithConfirm(filepath.Join(p.AbsPath, wr.path), []byte(wr.data)); err != nil {
			return simplerr.Wrap(err, "failed to write file")
		}
	}

	return nil
}
