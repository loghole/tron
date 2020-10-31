package helpers

import (
	"os/exec"

	"github.com/lissteron/simplerr"
)

func Exec(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "cmd '%s' failed", cmd.String())
	}

	return nil
}
