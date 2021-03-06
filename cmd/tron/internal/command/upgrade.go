package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/check"
	"github.com/loghole/tron/cmd/tron/internal/download"
	"github.com/loghole/tron/cmd/tron/internal/helpers"
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
		Use:   "upgrade [v1.2.3]",
		Short: "Self-upgrade tron tool",
		Example: "# upgrade by version tag:\n" +
			"tron upgrade v0.4.0\n" +
			"# upgrade to latest version:\n" +
			"tron upgrade\n" +
			"# get versions list:\n" +
			"tron upgrade --list",
		Run: u.run,
	}

	cmd.Flags().Bool(FlagList, false, "list available versions")
	cmd.Flags().Bool(FlagUnstable, false, "unstable tron versions <v1.2.3-rc2.0>")

	return cmd
}

func (u *Upgrade) run(cmd *cobra.Command, args []string) {
	unstable, err := cmd.Flags().GetBool(FlagUnstable)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	downloader, err := download.NewTron(u.printer, !unstable)
	if err != nil {
		u.printer.Printf(color.FgRed, "Create upgrader failed: %v\n", err)
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	version := "latest"

	if len(args) > 0 {
		version = args[0]
	}

	list, err := cmd.Flags().GetBool(FlagList)
	if err != nil {
		helpers.PrintCommandHelp(cmd)
		os.Exit(1)
	}

	if list {
		if err := downloader.ListVersions(); err != nil {
			u.printer.Printf(color.FgRed, "List versions failed: %v\n", err)
			os.Exit(1)
		}

		return
	}

	if ok := check.NewChecker(u.printer).CheckGolang(); !ok {
		u.printer.Println(color.FgRed, "Requirements check failed")
		os.Exit(1)
	}

	if err := downloader.Upgrade(version); err != nil {
		u.printer.Printf(color.FgRed, "Upgrade failed: %v\n", err)
		os.Exit(1)
	}

	u.printer.Println(color.FgGreen, "Success")
}
