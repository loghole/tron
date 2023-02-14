package generate

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Dockerfile(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate Dockerfile")

	path := filepath.Join(p.AbsPath, models.DockerfileFilepath)

	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerfileTemplate, p)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(dockerfile)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func DevDockerfile(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate Dockerfile.dev")

	path := filepath.Join(p.AbsPath, models.DockerfileDevFilepath)

	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerfileDev, p)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(dockerfile)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func DockerCompose(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate docker-compose.yaml")

	path := filepath.Join(p.AbsPath, models.DockerComposeFilepath)

	dockerfile, err := helpers.ExecTemplate(templates.DefaultDockerCompose, p)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(dockerfile)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func DockerComposeOverride(p *models.Project, printer stdout.Printer) error {
	printer.Println(color.FgMagenta, "Generate docker-compose.override.example.yaml")

	path := filepath.Join(p.AbsPath, models.DockerComposeOverrideFilepath)

	dockerfile, err := helpers.ExecTemplate(templates.DefaultOverrideDockerCompose, p)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if !helpers.ConfirmOverwrite(path) {
		printer.Println(color.FgBlue, "\tSkipped")

		return nil
	}

	if err := helpers.WriteToFile(path, []byte(dockerfile)); err != nil {
		return fmt.Errorf("write file '%s': %w", path, err)
	}

	printer.Println(color.FgBlue, "\tSuccess")

	return nil
}
