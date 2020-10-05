package generate

import (
	"bufio"
	"go/format"
	"os"
	"strings"

	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

type Config struct {
	imports map[string]struct{}
	files   []string
}

func NewConfig() *Config {
	return &Config{
		imports: make(map[string]struct{}),
		files: []string{
			// TODO: to const
			".deploy/config/values.yaml",
			".deploy/config/values_dev.yaml",
			".deploy/config/values_local.yaml",
			".deploy/config/values_stg.yaml",
			".deploy/config/values_prod.yaml",
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

	config, err := helpers.ExecTemplate(templates.ConfigTemplate, data)
	if err != nil {
		return simplerr.Wrap(err, "failed to exec template")
	}

	formatted, err := format.Source([]byte(config))
	if err != nil {
		return simplerr.Wrap(err, "failed to format config")
	}

	return helpers.WriteToFile("config/config.go", formatted)
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

		if strings.HasPrefix(parts[0], "#") {
			continue
		}

		c.imports[strings.TrimSpace(parts[0])] = struct{}{}
	}

	return nil
}
