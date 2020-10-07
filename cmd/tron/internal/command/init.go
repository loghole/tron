package command

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/generate"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	FlagProtoDirs = "proto"
	FlagConfig    = "config"
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
		Run:     i.run,
	}

	cmd.Flags().StringArray(FlagProtoDirs, []string{}, "directory with protos for generating your services")

	return cmd
}

func (i *InitCMD) run(cmd *cobra.Command, args []string) {
	if !project.NewChecker(i.printer).CheckRequirements() {
		panic("check requirements failed")
	}

	var (
		module string
		err    error
	)

	if len(args) > 0 {
		module = args[0]
	}

	protoDirs, err := cmd.Flags().GetStringArray(FlagProtoDirs)
	if err != nil {
		panic(err)
	}

	i.printer.VerbosePrintln(color.FgBlack, "Init project")

	i.project, err = project.NewProject(module)
	if err != nil {
		panic(err)
	}

	i.printer.VerbosePrintln(color.FgBlack, "Find proto files")

	if err := i.project.FindProtoFiles(protoDirs...); err != nil {
		panic(err)
	}

	i.printer.VerbosePrintln(color.FgBlack, "Move proto files")

	if err := i.project.MoveProtoFiles(); err != nil {
		panic(err)
	}

	if err := i.generate(generate.GoMod,
		generate.Makefile,
		generate.Linter,
		generate.Gitignore,
		generate.Dockerfile,
		generate.Values); err != nil {
		panic(err)
	}

	if err := helpers.Exec("make", "generate"); err != nil {
		panic(err)
	}

	if err := i.generate(generate.Config, generate.Mainfile); err != nil {
		panic(err)
	}
}

func (i *InitCMD) generate(list ...generate.Generator) error {
	for _, gen := range list {
		if err := gen(i.project, i.printer); err != nil {
			return err
		}
	}

	return nil
}
