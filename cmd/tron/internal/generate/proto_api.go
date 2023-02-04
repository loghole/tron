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

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

func ProtoAPI(project *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate proto api")

	// Buf lint.
	if err := lintProtos(project, printer); err != nil {
		return err
	}

	// Buf generate.
	if err := generateProtos(project, printer); err != nil {
		return err
	}

	// Move generated files.
	if err := moveGeneratedFiles(project, printer); err != nil {
		return err
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func lintProtos(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintf(color.Reset, "\trun buf lint: ")

	if err := runBufCommand(project, "lint"); err != nil {
		printer.VerbosePrintln(color.FgRed, "FAIL")
		printer.Println(color.FgRed, err)

		return ErrBufError
	}

	printer.VerbosePrintln(color.FgGreen, "OK")

	return nil
}

func generateProtos(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintf(color.Reset, "\trun buf generate: ")

	if err := runBufCommand(project, "generate"); err != nil {
		printer.VerbosePrintln(color.FgRed, "FAIL")
		printer.Println(color.FgRed, err)

		return ErrBufError
	}

	printer.VerbosePrintln(color.FgGreen, "OK")

	return nil
}

func runBufCommand(project *models.Project, args ...string) error {
	for _, file := range project.ProtoFiles {
		path := filepath.Join(
			models.ProjectPathVendorPB,
			project.Name,
			strings.TrimPrefix(file, project.AbsPath),
		)

		args = append(args, "--path", path)
	}

	output := bytes.NewBuffer(nil)

	cmd := exec.Command("./bin/buf", args...)
	cmd.Dir = project.AbsPath
	cmd.Stderr = output
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		out := strings.ReplaceAll(output.String(), models.ProjectPathVendorPB, "")
		out = strings.ReplaceAll("\t"+out, "\n", "\n\t")

		return fmt.Errorf(out) //nolint:goerr113,gocritic // need dynamic error for beautiful output
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
