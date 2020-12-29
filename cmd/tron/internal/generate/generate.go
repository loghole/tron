package generate

import (
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Generator func(p *models.Project, printer stdout.Printer) error
