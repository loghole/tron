package generate

import (
	"os"
	"path/filepath"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func GoMod(p *project.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.GoModFile)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return helpers.Exec("go", "mod", "init", p.Module)
		}

		return err
	}

	return nil
}
