package generate

import (
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Generator func(p *project.Project, printer stdout.Printer) error
