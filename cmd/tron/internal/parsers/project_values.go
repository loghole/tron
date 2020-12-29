package parsers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"

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

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")

		if len(parts) <= 1 || strings.HasPrefix(parts[0], "#") {
			continue
		}

		if _, ok := p.internal[strings.ToLower(parts[0])]; ok {
			continue
		}

		p.values[parts[0]] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
