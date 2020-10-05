package command

import (
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

	i.project, err = project.NewProject(module)
	if err != nil {
		panic(err)
	}

	if err := i.project.FindProtoFiles(protoDirs...); err != nil {
		panic(err)
	}

	if err := i.project.MoveProtoFiles(); err != nil {
		panic(err)
	}

	if err := i.project.InitGoMod(); err != nil {
		panic(err)
	}

	if err := i.project.InitMakeFile(); err != nil {
		panic(err)
	}

	if err := i.project.InitGitignore(); err != nil {
		panic(err)
	}

	if err := i.project.InitDockerfile(); err != nil {
		panic(err)
	}

	if err := i.project.InitValues(); err != nil {
		panic(err)
	}

	if err := helpers.Exec("make", "generate"); err != nil {
		panic(err)
	}

	if err := generate.NewConfig().Generate(); err != nil {
		panic(err)
	}

	if err := i.project.InitMainFile(); err != nil {
		panic(err)
	}
}
