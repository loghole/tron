package generate

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func ProtoAPI(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate proto api")

	// Buf generate.
	printer.VerbosePrintf(color.Reset, "\trun buf generate: ")

	if err := generateProtos(project); err != nil {
		printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

		return err
	}

	printer.VerbosePrintln(color.FgGreen, "OK")

	if err := moveGeneratedFiles(project, printer); err != nil {
		printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

		return err
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func generateProtos(project *models.Project) error {
	args := []string{
		"./bin/buf",
		"generate",
	}

	for _, file := range project.ProtoFiles {
		path := filepath.Join(
			models.ProjectPathVendorPB,
			project.Name,
			strings.TrimPrefix(file, project.AbsPath),
		)

		args = append(args, "--path", path)
	}

	stderr := bytes.NewBuffer(nil)

	cmd := exec.Command("./bin/buf", args...)
	cmd.Dir = project.AbsPath
	cmd.Stderr = stderr

	o, err := cmd.Output()
	if err != nil {
		return simplerr.Wrapf(err, "stderr: %s", stderr.String())
	}

	if len(o) > 0 {
		return simplerr.Wrapf(ErrBufUnexpectedOutput, string(o))
	}

	return nil
}

func moveGeneratedFiles(project *models.Project, printer stdout.Printer) error {
	path := filepath.Join(project.AbsPath, project.Name)

	defer os.RemoveAll(path)

	mover := func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		newPath := strings.TrimPrefix(path, project.AbsPath)
		newPath = strings.Replace(newPath, project.Name, models.ProjectPathPkgClients, 1)
		newPath = filepath.Join(project.AbsPath, newPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file '%s': %w", path, err)
		}

		if err := helpers.WriteToFile(newPath, data); err != nil {
			return fmt.Errorf("write file to '%s': %w", newPath, err)
		}

		printer.VerbosePrintf(color.Reset, "\tmove '%s' > '%s'\n", path, newPath)

		return nil
	}

	if err := filepath.Walk(path, mover); err != nil {
		return err
	}

	return nil
}
