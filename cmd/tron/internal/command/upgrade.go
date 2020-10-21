package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Upgrade struct {
	printer stdout.Printer
}

func NewUpgradeCMD(printer stdout.Printer) *Upgrade {
	return &Upgrade{printer: printer}
}

func (u *Upgrade) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Self-upgrade tron tool",
		Run:   u.run,
	}

	cmd.Flags().String(FlagVersion, "latest", "semver tag <v1.2.3>")

	return cmd
}

func (u *Upgrade) run(cmd *cobra.Command, args []string) {
	if ok := project.NewChecker(u.printer).CheckRequirements(); !ok {
		u.printer.Println(color.FgRed, "Requirements check failed")
		os.Exit(1)
	}

	panic("unimplemented")
}
