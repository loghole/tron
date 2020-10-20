package command

import (
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Upgrade struct {
	printer stdout.Printer
}

func NewUpgradeCMD(printer stdout.Printer) *Upgrade {
	return &Upgrade{printer: printer}
}

func (c *Upgrade) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Self-upgrade tron tool",
		Run:   c.run,
	}

	cmd.Flags().String(FlagVersion, "latest", "semver tag <v1.2.3>")

	return cmd
}

func (c *Upgrade) run(cmd *cobra.Command, args []string) {
	panic("unimplemented")
}
