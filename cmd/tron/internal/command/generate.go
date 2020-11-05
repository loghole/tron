package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/generate"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
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
		Short: "Generate project pkg and implementation from proto api",
		Long:  "Generate project pkg and implementation from proto api",
		Example: "# generate config constants from deploy values:\n" +
			"tron generate --config\n" +
			"# generate proto pkg and service implementations from protos:\n" +
			"tron generate --proto=api",
		Run: g.run,
	}

	cmd.Flags().StringArray(FlagProtoDirs, []string{}, "directory with protos for generating your services")
	cmd.Flags().Bool(FlagConfig, false, "Generate config helpers from values")

	return cmd
}

func (g *GenerateCMD) run(cmd *cobra.Command, args []string) {
	protoDirs, err := cmd.Flags().GetStringArray(FlagProtoDirs)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	g.project, err = project.NewProject("", g.printer)
	if err != nil {
		g.printer.Printf(color.FgRed, "Parse project failed: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	if len(protoDirs) > 0 {
		if err := g.runProto(protoDirs); err != nil {
			g.printer.Printf(color.FgRed, "Generate protos failed: %v\n", err)
			helpers.PrintCommandHelp(cmd)
			os.Exit(1)
		}
	}

	config, err := cmd.Flags().GetBool(FlagConfig)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	if config {
		if err := g.runConfig(); err != nil {
			g.printer.Printf(color.FgRed, "Generate config failed: %v\n", err)
			helpers.PrintCommandHelp(cmd)
			os.Exit(1)
		}
	}

	g.printer.Println(color.FgGreen, "Success")
}

func (g *GenerateCMD) runProto(dirs []string) (err error) {
	if err := g.project.FindProtoFiles(dirs...); err != nil {
		return simplerr.Wrap(err, "find proto files failed")
	}

	if g.project.WithoutProtos() {
		return nil
	}

	if ok := project.NewChecker(g.printer).CheckProtoc(); !ok {
		g.printer.Println(color.FgRed, "Requirements check failed")
		os.Exit(1)
	}

	if err := generate.VendorPB(g.project, g.printer); err != nil {
		return simplerr.Wrap(err, "download proto imports failed")
	}

	if err := generate.Protos(g.project, g.printer); err != nil {
		return simplerr.Wrap(err, "generate proto files failed")
	}

	return nil
}

func (g *GenerateCMD) runConfig() error {
	if err := generate.Config(g.project, g.printer); err != nil {
		return simplerr.Wrap(err, "generate config failed")
	}

	g.printer.VerbosePrintln(color.FgCyan, "\tCreated")

	return nil
}
