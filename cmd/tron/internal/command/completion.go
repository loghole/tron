package command

import (
	"os"

	"github.com/spf13/cobra"
)

type CompletionCMD struct{}

func NewCompletionCMD() *CompletionCMD {
	return &CompletionCMD{}
}

func (c *CompletionCMD) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <(tron completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ tron completion bash > /etc/bash_completion.d/tron
  # macOS:
  $ tron completion bash > /usr/local/etc/bash_completion.d/tron

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ tron completion zsh > ~/.oh-my-zsh/completions/_tron

  # You will need to start a new shell for this setup to take effect.

PowerShell:

  PS> tron completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> tron completion powershell > tron.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run:                   c.run,
	}

	return cmd
}

func (c *CompletionCMD) run(cmd *cobra.Command, args []string) {
	switch args[0] {
	case "bash":
		_ = cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		_ = cmd.Root().GenZshCompletion(os.Stdout)
	case "powershell":
		_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	}
}
