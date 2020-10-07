package generate

import (
	"strings"

	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
	"github.com/loghole/tron/internal/app"
)

func Values(p *project.Project, printer stdout.Printer) error {
	data := templates.ValuesData{
		List: []templates.Env{
			{Key: strings.ToUpper(app.LoggerLevelEnv), Val: "info"},
			{Key: strings.ToUpper(app.LoggerCollectorAddrEnv), Val: ""},
			{Key: strings.ToUpper(app.LoggerDisableStdoutEnv), Val: "false"},
			{Key: strings.ToUpper(app.JaegerAddrEnv), Val: "127.0.0.1:6831"},
			{Key: strings.ToUpper(app.HTTPPortEnv), Val: "8080"},
			{Key: strings.ToUpper(app.GRPCPortEnv), Val: "8081"},
			{Key: strings.ToUpper(app.AdminPortEnv), Val: "8082"},
		},
	}

	values, err := helpers.ExecTemplate(templates.ValuesTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	if err := helpers.WriteWithConfirm(models.ValuesBaseFilepath, []byte(values)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteWithConfirm(models.ValuesDevFilepath, []byte(templates.ValuesDevTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteWithConfirm(models.ValuesLocalFilepath, []byte(templates.ValuesLocalTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteWithConfirm(models.ValuesStgFilepath, []byte(templates.ValuesStgTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	if err := helpers.WriteWithConfirm(models.ValuesProdFilepath, []byte(templates.ValuesProdTemplate)); err != nil {
		return simplerr.Wrap(err, "failed to write file")
	}

	return nil
}
