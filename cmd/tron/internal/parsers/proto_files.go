package parsers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type ProtoFilesParser struct {
	project *models.Project
	printer stdout.Printer
}

func NewProtoFilesParser(project *models.Project, printer stdout.Printer) *ProtoFilesParser {
	return &ProtoFilesParser{
		project: project,
		printer: printer,
	}
}

func (p *ProtoFilesParser) Parse() error {
	p.printer.VerbosePrintln(color.FgMagenta, "Find proto files")

	for _, dir := range p.project.ProtoDirs {
		absPath := filepath.Join(p.project.AbsPath, dir)

		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err
		}

		if err := filepath.Walk(absPath, p.walkFn); err != nil {
			return err
		}
	}

	p.printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func (p *ProtoFilesParser) walkFn(path string, info os.FileInfo, err error) error {
	switch {
	case err != nil:
		return err
	case info.IsDir():
		return nil
	case filepath.Ext(path) != models.ProtoExt:
		return nil
	case strings.Contains(path, models.ProjectPathVendorPB):
		return nil
	}

	p.project.ProtoFiles = append(p.project.ProtoFiles, path)

	p.printer.VerbosePrintf(color.Reset, "\tcollected proto '%s'\n", color.YellowString(path))

	return nil
}
