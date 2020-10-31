package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Check struct {
	printer stdout.Printer
}

func NewCheckCMD(printer stdout.Printer) *Check {
	return &Check{printer: printer}
}

func (c *Check) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Check system requirements",
		Long:    "Check if system compatible with current version.",
		Example: "tron check",
		Run:     c.run,
	}
}

func (c *Check) run(cmd *cobra.Command, args []string) {
	if ok := project.NewChecker(c.printer).CheckRequirements(); !ok {
		c.printer.Println(color.FgHiRed, "Requirements check failed")
		os.Exit(1)
	}

	c.printer.Println(color.FgGreen, "Success\n")
}
