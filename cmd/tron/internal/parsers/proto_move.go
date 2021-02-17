package parsers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type ProtoFilesMover struct {
	project *models.Project
	printer stdout.Printer
}

func NewProtoFilesMover(project *models.Project, printer stdout.Printer) *ProtoFilesMover {
	return &ProtoFilesMover{
		project: project,
		printer: printer,
	}
}

func (m *ProtoFilesMover) Move() error {
	m.printer.VerbosePrintln(color.FgMagenta, "Move proto files")

	for idx, path := range m.project.ProtoFiles {
		if strings.Contains(path, models.ProjectPathVendorPB) {
			continue
		}

		version, err := m.findPackageVersion(path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(m.project.AbsPath, models.ProjectPathAPI, version, filepath.Base(path))

		if path == newPath {
			continue
		}

		m.printer.VerbosePrintf(color.Reset, "\tmove '%s' > '%s'", path, newPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := helpers.WriteToFile(newPath, data); err != nil {
			return err
		}

		if err := os.Remove(path); err != nil {
			return err
		}

		m.project.ProtoFiles[idx] = newPath
	}

	m.project.ProtoDirs = []string{models.ProjectPathAPI}

	m.printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func (m *ProtoFilesMover) findPackageVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file '%s': %w", path, err)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := models.PackageRegexp.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
		if len(m) <= 1 {
			continue
		}

		m = models.ProtoPkgVersionRegexp.FindStringSubmatch(m[1])
		if len(m) <= 1 {
			return "", ErrInvalidProtoPkgName
		}

		return m[1], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", ErrPackageNotPound
}
