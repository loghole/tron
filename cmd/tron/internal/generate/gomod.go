package generate

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func GoMod(p *project.Project, printer stdout.Printer) error {
	path := filepath.Join(p.AbsPath, models.GoModFile)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			args := []string{"mod", "init", p.Module}

			return exec.Command("go", args...).Run()
		}

		return err
	}

	return nil
}
