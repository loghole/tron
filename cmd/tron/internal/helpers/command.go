package helpers

import (
	"github.com/spf13/cobra"
)

func PrintCommandHelp(cmd *cobra.Command) {
	_ = cmd.Help()
}
