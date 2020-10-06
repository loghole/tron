package generate

import (
	"bufio"
	"go/format"
	"os"
	"strings"

	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/templates"
	"github.com/loghole/tron/internal/app"
)

type Config struct {
	imports  map[string]struct{}
	internal map[string]struct{}
	files    []string
}

func NewConfig() *Config {
	return &Config{
		imports: make(map[string]struct{}),
		files: []string{
			models.ValuesBaseFilepath,
			models.ValuesDevFilepath,
			models.ValuesLocalFilepath,
			models.ValuesStgFilepath,
			models.ValuesProdFilepath,
		},
		internal: map[string]struct{}{
			app.NamespaceEnv: {},
			app.LoggerLevelEnv: {},
			app.LoggerCollectorAddrEnv: {},
			app.LoggerDisableStdoutEnv: {},
			app.JaegerAddrEnv: {},
			app.HTTPPortEnv: {},
			app.GRPCPortEnv: {},
			app.AdminPortEnv: {},
		},
	}
}

func (c *Config) Generate() error {
	for _, file := range c.files {
		if err := c.parseFile(file); err != nil {
			return err
		}
	}

	data := &templates.ConfigData{Values: make([]templates.ConfigValue, 0, len(c.imports))}

	for key := range c.imports {
		data.Values = append(data.Values, templates.NewConfigValue(key))
	}

	config, err := helpers.ExecTemplate(templates.ConfigConstTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formatted, err := format.Source([]byte(config))
	if err != nil {
		return simplerr.Wrap(err, "failed to format config")
	}

	if err := helpers.WriteToFile(models.ConfigConstFilepath, formatted); err != nil {
		return err
	}

	if err := helpers.WriteToFile(models.ConfigFilepath, []byte(templates.ConfigTemplate)); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return simplerr.Wrapf(err, "open file '%s' failed", filepath)
	}

	defer helpers.Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		parts := strings.Split(scanner.Text(), ": ")

		if len(parts) <= 1 {
			continue
		}

		name := strings.TrimSpace(parts[0])

		if strings.HasPrefix(name, "#") {
			continue
		}

		if _, ok := c.internal[name]; ok {
			continue
		}

		c.imports[name] = struct{}{}
	}

	return nil
}
