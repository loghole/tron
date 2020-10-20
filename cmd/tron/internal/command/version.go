package command

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

// lint:gochecknoglobals //tron version
var version = "+devel"

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
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			v.printer.Printf(color.FgBlack, "Tron version: %s\n", color.YellowString(version))
		},
	}
}
