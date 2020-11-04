package helpers

import (
	"os"
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

func ExecWithPrint(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "cmd '%s' failed", cmd.String())
	}

	return nil
}
