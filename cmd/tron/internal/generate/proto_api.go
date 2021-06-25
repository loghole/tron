package generate

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

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

	// Move pkg clients.
	printer.VerbosePrintf(color.Reset, "\tmv generated files: ")

	if err := moveGeneratedFile(project); err != nil {
		printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

		return err
	}

	printer.VerbosePrintln(color.FgGreen, "OK")

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

	cmd := exec.Command(args[0], args[1:]...)
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

func moveGeneratedFile(project *models.Project) error {
	return os.Rename(
		filepath.Join(project.AbsPath, project.Name),
		filepath.Join(project.AbsPath, models.ProjectPathPkgClients),
	)
}
