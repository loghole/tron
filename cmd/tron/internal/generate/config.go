package generate

import (
	"bufio"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
	"github.com/loghole/tron/cmd/tron/internal/version"
	"github.com/loghole/tron/internal/app"
)

type config struct {
	project  *project.Project
	printer  stdout.Printer
	imports  map[string]struct{}
	internal map[string]struct{}
	files    []string
}

func Config(p *project.Project, printer stdout.Printer) error {
	generator := &config{
		project: p,
		printer: printer,
		imports: make(map[string]struct{}),
		files: []string{
			models.ValuesBaseFilepath,
			models.ValuesDevFilepath,
			models.ValuesLocalFilepath,
			models.ValuesStgFilepath,
			models.ValuesProdFilepath,
		},
		internal: map[string]struct{}{
			app.NamespaceEnv:           {},
			app.LoggerLevelEnv:         {},
			app.LoggerCollectorAddrEnv: {},
			app.LoggerDisableStdoutEnv: {},
			app.JaegerAddrEnv:          {},
			app.HTTPPortEnv:            {},
			app.GRPCPortEnv:            {},
			app.AdminPortEnv:           {},
		},
	}

	return generator.run()
}

func (c *config) run() error {
	c.printer.VerbosePrintln(color.FgMagenta, "Generate config")

	for _, file := range c.files {
		if err := c.parseFile(file); err != nil {
			return err
		}
	}

	data := &templates.ConfigData{
		TronVersion: version.CliVersion,
		Values:      make([]templates.ConfigValue, 0, len(c.imports)),
	}

	for key := range c.imports {
		data.Values = append(data.Values, templates.NewConfigValue(key))
	}

	data.SortValues()

	config, err := helpers.ExecTemplate(templates.ConfigConstTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formatted, err := format.Source([]byte(config))
	if err != nil {
		return simplerr.Wrap(err, "failed to format config")
	}

	if err := helpers.WriteToFile(
		filepath.Join(c.project.AbsPath, models.ConfigConstFilepath),
		formatted); err != nil {
		return err
	}

	if err := helpers.WriteWithConfirm(
		filepath.Join(c.project.AbsPath, models.ConfigFilepath),
		[]byte(templates.ConfigTemplate)); err != nil {
		return err
	}

	return nil
}

func (c *config) parseFile(path string) error {
	file, err := os.Open(filepath.Join(c.project.AbsPath, path))
	if err != nil {
		return simplerr.Wrapf(err, "open file '%s' failed", path)
	}

	defer helpers.Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		parts := strings.Split(scanner.Text(), ":")

		if len(parts) <= 1 {
			continue
		}

		name := strings.TrimSpace(parts[0])

		if strings.HasPrefix(name, "#") {
			continue
		}

		if _, ok := c.internal[strings.ToLower(name)]; ok {
			continue
		}

		c.imports[name] = struct{}{}
	}

	return nil
}
