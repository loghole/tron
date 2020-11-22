package command

import (
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/generate"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	FlagProtoDirs = "proto"
	FlagConfig    = "config"
	FlagVersion   = "version"
	FlagList      = "list"
	FlagUnstable  = "unstable"
)

type InitCMD struct {
	printer stdout.Printer
	project *project.Project
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
	if ok := project.NewChecker(i.printer).CheckInitRequirements(); !ok {
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

	if err := i.runInit(module, protoDirs); err != nil {
		i.printer.Printf(color.FgRed, "Init failed: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	i.printer.Println(color.FgGreen, "Success")
}

func (i *InitCMD) runInit(module string, dirs []string) (err error) {
	i.printer.VerbosePrintln(color.FgMagenta, "Init project")

	i.project, err = project.NewProject(module, i.printer)
	if err != nil {
		return simplerr.Wrap(err, "parse project failed")
	}

	if err := i.project.FindProtoFiles(dirs...); err != nil {
		return simplerr.Wrap(err, "find proto files failed")
	}

	if err := i.project.MoveProtoFiles(); err != nil {
		return simplerr.Wrap(err, "move proto files failed")
	}

	if err := i.generate(
		generate.Git,
		generate.GoMod,
		generate.TronMK,
		generate.Makefile,
		generate.Linter,
		generate.Gitignore,
		generate.Dockerfile,
		generate.Values,
		generate.ReadmeMD); err != nil {
		return simplerr.Wrap(err, "generate files failed")
	}

	i.printer.Println(color.FgMagenta, "Generate files from proto api if exists")

	if err := helpers.ExecWithPrint(i.project.AbsPath, "make", "generate"); err != nil {
		return simplerr.Wrap(err, "exec 'make generate' failed")
	}

	i.printer.Println(color.FgMagenta, "Generate config and main files")

	if err := i.generate(
		generate.Config,
		generate.ConfigHelper,
		generate.Mainfile); err != nil {
		return simplerr.Wrap(err, "generate config and main files failed")
	}

	if err := helpers.Exec(i.project.AbsPath, "go", "mod", "tidy"); err != nil {
		return simplerr.Wrap(err, "exec 'go mod tidy' failed")
	}

	return nil
}

func (i *InitCMD) generate(list ...generate.Generator) error {
	for _, gen := range list {
		if err := gen(i.project, i.printer); err != nil {
			if errors.Is(err, generate.ErrAlreadyExists) {
				continue
			}

			return err
		}

		i.printer.VerbosePrintln(color.FgCyan, "\tCreated")
	}

	return nil
}
