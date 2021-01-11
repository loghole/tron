package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/download"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

type VersionCMD struct {
	printer stdout.Printer
}

func NewVersionCMD(printer stdout.Printer) *VersionCMD {
	return &VersionCMD{printer: printer}
}

func (v *VersionCMD) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print tron version",
		Example: "tron version",
		Run:     v.run,
	}

	cmd.Flags().Bool(FlagList, false, "list available versions")

	return cmd
}

func (v *VersionCMD) run(cmd *cobra.Command, args []string) {
	list, err := cmd.Flags().GetBool(FlagList)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	if list {
		downloader, err := download.NewTron(v.printer, true)
		if err != nil {
			v.printer.Printf(color.FgRed, "Create upgrader failed: %v\n", err)
			helpers.PrintCommandHelp(cmd)
			os.Exit(1)
		}

		if err := downloader.ListVersions(); err != nil {
			v.printer.Printf(color.FgRed, "List versions failed: %v\n", err)
			os.Exit(1)
		}

		return
	}

	v.printer.Printf(color.Reset, "Tron version: %s\n", color.YellowString(version.CliVersion))
}
