package command

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/check"
	"github.com/loghole/tron/cmd/tron/internal/download"
	"github.com/loghole/tron/cmd/tron/internal/generate"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/parsers"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type InitCMD struct {
	printer stdout.Printer
}

func NewInitCMD(printer stdout.Printer) *InitCMD {
	return &InitCMD{printer: printer}
}

func (i *InitCMD) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [module]",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize Application",
		Long:    "Initialize will create a new application",
		Example: "# from project dir with proto files:\n" +
			"tron init github.com/loghole/example --proto api\n" +
			"# from root dir with create project dir:\n" +
			"tron init github.com/loghole/example",
		Run: i.run,
	}

	cmd.Flags().StringArray(FlagProtoDirs, []string{}, "directory with protos for generating your services")

	return cmd
}

func (i *InitCMD) run(cmd *cobra.Command, args []string) {
	if ok := check.NewChecker(i.printer).CheckInitRequirements(); !ok {
		i.printer.Println(color.FgRed, "Requirements check failed")
		os.Exit(1)
	}

	var module string

	if len(args) > 0 {
		module = args[0]
	}

	protoDirs, err := cmd.Flags().GetStringArray(FlagProtoDirs)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	if err := i.runInitService(parsers.WithModuleName(module), parsers.WithProtoDirs(protoDirs)); err != nil {
		i.printer.Printf(color.FgRed, "Init failed: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	i.printer.Println(color.FgGreen, "Success")
}

func (i *InitCMD) runInitService(opts ...parsers.Option) (err error) {
	i.printer.VerbosePrintln(color.FgMagenta, "Init project")

	project, err := parsers.NewProjectParser(i.printer, opts...).Parse()
	if err != nil {
		return fmt.Errorf("parse project: %w", err)
	}

	if err := parsers.NewProtoFilesMover(project, i.printer).Move(); err != nil {
		return fmt.Errorf("move proto files: %w", err)
	}

	if err := download.NewVendor(project, i.printer).Download(); err != nil {
		return fmt.Errorf("vendor proto files: %w", err)
	}

	if err := i.generate(
		project,
		generate.Git,
		generate.Buf,
		generate.GoMod,
		generate.TronMK,
		generate.Makefile,
		generate.Linter,
		generate.Gitignore,
		generate.Dockerfile,
		generate.DevDockerfile,
		generate.DockerCompose,
		generate.DockerComposeOverride,
		generate.ReadmeMD,
		generate.MainFile,
	); err != nil {
		return fmt.Errorf("generate files: %w", err)
	}

	i.printer.Println(color.FgMagenta, "Generate files from proto api if exists")

	if err := helpers.ExecWithPrint(project.AbsPath,
		"make",
		"docker-compose",
		"docker-volumes",
		"tidy",
		"generate",
	); err != nil {
		return fmt.Errorf("exec 'make generate': %w", err)
	}

	return nil
}

func (i *InitCMD) generate(project *models.Project, list ...generate.Generator) error {
	for _, gen := range list {
		if err := gen(project, i.printer); err != nil {
			return err
		}
	}

	return nil
}
