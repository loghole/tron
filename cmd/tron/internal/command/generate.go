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

type GenerateCMD struct {
	printer stdout.Printer
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
	cmd.Flags().Bool(FlagImpl, true, "Generate grpc api implementation")

	return cmd
}

func (g *GenerateCMD) run(cmd *cobra.Command, args []string) {
	// Parse flags.
	protoDirs, err := cmd.Flags().GetStringArray(FlagProtoDirs)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	config, err := cmd.Flags().GetBool(FlagConfig)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	impl, err := cmd.Flags().GetBool(FlagImpl)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	// Parse project.
	project, err := parsers.NewProjectParser(g.printer, parsers.WithProtoDirs(protoDirs)).Parse()
	if err != nil {
		g.printer.Printf(color.FgRed, "Parse project failed: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	// Generate.
	if config {
		if err := g.runConfig(project); err != nil {
			g.printer.Printf(color.FgRed, "Generate config failed: %v\n", err)
			helpers.PrintCommandHelp(cmd)
			os.Exit(1)
		}
	}

	if len(protoDirs) > 0 && project.WithProtos() {
		if ok := check.NewChecker(g.printer).CheckProtoc(); !ok {
			g.printer.Println(color.FgRed, "Requirements check failed")
			os.Exit(1)
		}

		if err := g.runProto(project, impl); err != nil {
			g.printer.Printf(color.FgRed, "Generate protos failed: %v\n", err)
			helpers.PrintCommandHelp(cmd)
			os.Exit(1)
		}
	}

	if err := generate.TronMK(project, g.printer); err != nil {
		g.printer.Printf(color.FgRed, "Generate tron mk: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	g.printer.Println(color.FgGreen, "Success")
}

func (g *GenerateCMD) runProto(project *models.Project, impl bool) (err error) {
	if err := download.NewVendor(project, g.printer).Download(); err != nil {
		return fmt.Errorf("vandor proto files: %w", err)
	}

	if err := parsers.NewProtoDescParser(project, g.printer).Parse(); err != nil {
		return fmt.Errorf("parse proto files: %w", err)
	}

	if err := generate.ProtoAPI(project, g.printer); err != nil {
		return fmt.Errorf("generate proto api: %w", err)
	}

	if !impl {
		return nil
	}

	if err := generate.Implement(project, g.printer); err != nil {
		return fmt.Errorf("generate implement: %w", err)
	}

	return nil
}

func (g *GenerateCMD) runConfig(project *models.Project) error {
	if err := parsers.NewValuesParser(project, g.printer).Parser(); err != nil {
		return fmt.Errorf("parse config values: %w", err)
	}

	if err := generate.Config(project, g.printer); err != nil {
		return fmt.Errorf("generate config: %w", err)
	}

	return nil
}
