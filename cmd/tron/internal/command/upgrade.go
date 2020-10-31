package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/upgrade"
)

type Upgrade struct {
	printer stdout.Printer
}

func NewUpgradeCMD(printer stdout.Printer) *Upgrade {
	return &Upgrade{printer: printer}
}

func (u *Upgrade) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Self-upgrade tron tool",
		Example: u.example(),
		Run:     u.run,
	}

	cmd.Flags().String(FlagVersion, "latest", "semver tag <v1.2.3>")
	cmd.Flags().Bool(FlagList, false, "list available versions")

	return cmd
}

func (u *Upgrade) example() string {
	return "# upgrade by version tag:\n" +
		"tron upgrade --version v0.4.0\n" +
		"# upgrade to latest version:\n" +
		"tron upgrade\n"+
		"# get versions list:\n" +
		"tron upgrade --list"

}

func (u *Upgrade) run(cmd *cobra.Command, args []string) {
	upgrader, err := upgrade.New(u.printer)
	if err != nil {
		u.printer.Println(color.FgRed, "create upgrader failed")
		os.Exit(1)
	}

	version, err := cmd.Flags().GetString(FlagVersion)
	if err != nil {
		panic(err)
	}

	list, err := cmd.Flags().GetBool(FlagList)
	if err != nil {
		panic(err)
	}

	if list {
		if err := upgrader.ListVersions(); err != nil {
			u.printer.Println(color.FgRed, "list versions failed")
			os.Exit(1)
		}

		return
	}

	if ok := project.NewChecker(u.printer).CheckGolang(); !ok {
		u.printer.Println(color.FgRed, "Requirements check failed")
		os.Exit(1)
	}

	if err := upgrader.Upgrade(version); err != nil {
		u.printer.Printf(color.FgRed, "Upgrade failed: %v\n", err)
		os.Exit(1)
	}
}
