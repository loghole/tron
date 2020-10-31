package command

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

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
	return &cobra.Command{
		Use:     "version",
		Short:   "Print tron version",
		Example: "tron version",
		Run: func(cmd *cobra.Command, args []string) {
			v.printer.Printf(color.Reset, "Tron version: %s\n", color.YellowString(version.CliVersion))
		},
	}
}
