package parsers

import (
	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

type Option func(p *ProjectParser)

func WithModuleName(name string) Option {
	return func(p *ProjectParser) {
		p.project.Module = name
	}
}

func WithProtoDirs(dirs []string) Option {
	return func(p *ProjectParser) {
		p.project.ProtoDirs = dirs
	}
}

type ProjectParser struct {
	project *models.Project
	printer stdout.Printer
}

func NewProjectParser(printer stdout.Printer, opts ...Option) *ProjectParser {
	parser := &ProjectParser{
		printer: printer,
		project: &models.Project{
			Version: version.CliVersion,
		},
	}

	for _, opt := range opts {
		opt(parser)
	}

	return parser
}

func (p *ProjectParser) Parse() (*models.Project, error) {
	p.printer.VerbosePrintln(color.FgMagenta, "Parse project")

	if err := NewModuleParser(p.project, p.printer).Parse(); err != nil {
		return nil, err
	}

	if err := NewProtoFilesParser(p.project, p.printer).Parse(); err != nil {
		return nil, err
	}

	return p.project, nil
}
