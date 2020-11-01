package helpers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PrintCommandHelp(cmd *cobra.Command) {
	fmt.Println()
	_ = cmd.Help()
}
