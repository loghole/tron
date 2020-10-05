package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func ConfirmOverwrite(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return true
	}

	var (
		keep  = fmt.Sprintf("Keep current version of `%s`", filepath.Base(path))
		write = fmt.Sprintf("Write default version of `%s`", filepath.Base(path))
	)

	prompt := promptui.Select{
		Label: color.BlueString("'%s' already exists, do you want to keep this file", path),
		Items: []string{keep, write},
	}

	_, result, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			os.Exit(130)
		}

		color.Red("Prompt failed %v: %s", err, keep)

		return false
	}

	return result == write
}
