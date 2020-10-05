package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/generate"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type GenerateCMD struct {
	printer stdout.Printer
	project *project.Project
}

func NewGenerateCMD(printer stdout.Printer) *GenerateCMD {
	return &GenerateCMD{printer: printer}
}

func (g *GenerateCMD) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate...",
		Long:  "Various generate...",
		Run:   g.run,
	}

	cmd.Flags().StringArray(FlagProtoDirs, []string{}, "directory with protos for generating your services")

	return cmd
}

func (g *GenerateCMD) run(cmd *cobra.Command, args []string) {
	protoDirs, err := cmd.Flags().GetStringArray(FlagProtoDirs)
	if err != nil {
		panic(err)
	}

	if len(protoDirs) > 0 {
		if err := g.runProto(protoDirs); err != nil {
			color.Red("Generate protos failed: %v", err)
			os.Exit(1)
		}
	}
}

func (g *GenerateCMD) runProto(dirs []string) error {
	if ok := project.NewChecker(g.printer).CheckRequirements(); !ok {
		color.Red("\nRequirements check failed")
		os.Exit(1)
	}

	var err error

	g.printer.VerbosePrintln(color.FgBlack, "Parse project")

	g.project, err = project.NewProject("")
	if err != nil {
		color.Red("Parse project failed: %v", err)
		os.Exit(1)
	}

	g.printer.VerbosePrintln(color.FgBlack, "Find proto files")

	if err := g.project.FindProtoFiles(dirs...); err != nil {
		color.Red("Find proto files failed: %v", err)
		os.Exit(1)
	}

	g.printer.VerbosePrintln(color.FgBlack, "Start vendoring")

	if err := generate.NewVendorPB(g.project).Download(); err != nil {
		color.Red("Download proto imports failed: %v", err)
		os.Exit(1)
	}

	g.printer.VerbosePrintln(color.FgBlack, "Generate profiles")

	if err := generate.NewProto(g.project).Generate(); err != nil {
		color.Red("Generate proto files failed: %v", err)
		os.Exit(1)
	}

	return nil
}
