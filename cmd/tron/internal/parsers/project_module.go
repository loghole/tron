package parsers

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type ModuleParser struct {
	project *models.Project
	printer stdout.Printer
}

func NewModuleParser(project *models.Project, printer stdout.Printer) *ModuleParser {
	return &ModuleParser{
		project: project,
		printer: printer,
	}
}

func (p *ModuleParser) Parse() (err error) {
	if p.project.Module == "" {
		p.project.Module, err = p.parseGoMod()
		if err != nil {
			return err
		}
	}

	if p.project.Module == "" {
		return ErrModuleNotFound
	}

	absPath, err := os.Getwd()
	if err != nil {
		return err
	}

	if !strings.HasSuffix(absPath, helpers.ModuleName(p.project.Module)) {
		absPath = filepath.Join(absPath, helpers.ModuleName(p.project.Module))

		if err := helpers.MkdirWithConfirm(absPath); err != nil {
			return err
		}
	}

	p.project.Name = helpers.ModuleName(p.project.Module)
	p.project.AbsPath = absPath

	return nil
}

func (p *ModuleParser) parseGoMod() (string, error) {
	file, err := os.Open(models.GoModFile)
	if err != nil {
		return "", err
	}

	defer helpers.Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if m := models.ModuleRegexp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", ErrModuleNotFound
}
