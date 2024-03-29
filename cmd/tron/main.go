package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/check"
	"github.com/loghole/tron/cmd/tron/internal/command"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	FlagVerbose = "verbose"
)

func main() {
	if err := os.Setenv("GO111MODULE", "on"); err != nil {
		color.Red("Set GO111MODULE=on failed: %v", err)
		os.Exit(1)
	}

	var (
		printer       = stdout.NewPrinter()
		checkCMD      = command.NewCheckCMD(printer)
		initCMD       = command.NewInitCMD(printer)
		generateCMD   = command.NewGenerateCMD(printer)
		upgradeCMD    = command.NewUpgradeCMD(printer)
		versionCMD    = command.NewVersionCMD(printer)
		completionCMD = command.NewCompletionCMD()
	)

	check.NewChecker(printer).CheckTron()

	rootCmd := &cobra.Command{
		Use:   "tron",
		Short: "A generator for Tron based Applications",
		Long:  "Tron is a CLI library for generating GO services.",
	}

	rootCmd.AddCommand(
		checkCMD.Command(),
		initCMD.Command(),
		generateCMD.Command(),
		upgradeCMD.Command(),
		versionCMD.Command(),
		completionCMD.Command(),
	)

	rootCmd.PersistentFlags().BoolP(FlagVerbose, "v", false, "make tron more verbose")

	cobra.OnInitialize(func() {
		verbose, err := rootCmd.Flags().GetBool("verbose")
		if err != nil {
			printer.Printf(color.FgYellow, "can't get `verbose` flag: %s\n", err)
		}

		if verbose {
			printer.Verbose()
		}
	})

	if err := rootCmd.Execute(); err != nil {
		color.Red("Exec root command failed: %v", err)
		os.Exit(1)
	}
}
