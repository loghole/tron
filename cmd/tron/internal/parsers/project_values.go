package parsers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/internal/app"
)

type ValuesParser struct {
	project *models.Project
	printer stdout.Printer

	values   map[string]struct{}
	internal map[string]struct{}
	files    []string
}

func NewValuesParser(project *models.Project, printer stdout.Printer) *ValuesParser {
	return &ValuesParser{
		project: project,
		printer: printer,
		values:  make(map[string]struct{}),
		internal: map[string]struct{}{
			app.NamespaceEnv:           {},
			app.LoggerLevelEnv:         {},
			app.LoggerCollectorAddrEnv: {},
			app.LoggerDisableStdoutEnv: {},
			app.JaegerAddrEnv:          {},
			app.JaegerSamplerType:      {},
			app.JaegerSamplerParam:     {},
			app.HTTPPortEnv:            {},
			app.GRPCPortEnv:            {},
			app.AdminPortEnv:           {},
		},
		files: []string{
			models.ValuesBaseFilepath,
			models.ValuesDevFilepath,
			models.ValuesLocalFilepath,
			models.ValuesStgFilepath,
			models.ValuesProdFilepath,
		},
	}
}

func (p *ValuesParser) Parser() error {
	p.printer.VerbosePrintln(color.FgMagenta, "Parse config values")

	if err := p.parseFiles(); err != nil {
		return err
	}

	result := make([]string, 0, len(p.values))

	for val := range p.values {
		result = append(result, val)
	}

	sort.Strings(result)

	for _, val := range result {
		p.project.ValuesEnv = append(p.project.ValuesEnv, &models.ConfigValue{
			Name: helpers.UpperCamelCase(helpers.GoName(val)),
			Key:  strings.ToLower(val),
		})
	}

	p.printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func (p *ValuesParser) parseFiles() error {
	for _, path := range p.files {
		if err := p.parseFile(path); err != nil {
			return fmt.Errorf("parse file '%s': %w", path, err)
		}
	}

	return nil
}

func (p *ValuesParser) parseFile(path string) error {
	file, err := os.Open(filepath.Join(p.project.AbsPath, path))
	if err != nil {
		return fmt.Errorf("open file '%s': %w", path, err)
	}

	defer helpers.Close(file)

	var dest map[string]interface{}

	if err := yaml.NewDecoder(file).Decode(&dest); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return fmt.Errorf("decode yaml config: %w", err)
	}

	for key := range dest {
		if _, ok := p.internal[strings.ToUpper(key)]; ok {
			continue
		}

		p.values[key] = struct{}{}
	}

	return nil
}
