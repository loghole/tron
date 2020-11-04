package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"github.com/manifoldco/promptui"
)

var (
	ErrIsNotDir     = errors.New("is not dir")
	ErrUnknownInput = errors.New("unknown input")
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
			os.Exit(1)
		}

		color.Red("Prompt failed %v: %s", err, keep)

		return false
	}

	return result == write
}

func WriteWithConfirm(path string, data []byte) error {
	if !ConfirmOverwrite(path) {
		return nil
	}

	return WriteToFile(path, data)
}

func MkdirWithConfirm(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if os.IsNotExist(err) {
		return Mkdir(path + string(filepath.Separator))
	}

	if !stat.IsDir() {
		return simplerr.Wrap(ErrIsNotDir, path)
	}

	var (
		cancel    = "Cancel"
		overwrite = "Overwrite dir"
		merge     = "Merge dirs"
	)

	prompt := promptui.Select{
		Label: color.BlueString("'%s' already exists", path),
		Items: []string{cancel, overwrite, merge},
	}

	_, result, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			os.Exit(1)
		}

		return err
	}

	switch result {
	case cancel:
		os.Exit(1)
	case overwrite:
		if err := os.RemoveAll(path); err != nil {
			return err
		}

		return Mkdir(path + string(filepath.Separator))
	case merge:
		return nil
	}

	return ErrUnknownInput
}
