package allowlist

import (
	"tynmo/command/allowlist/deployment"
	"tynmo/command/allowlist/show"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	allowlistCmd := &cobra.Command{
		Use:   "allowlist",
		Short: "Top level command for modifying the Tynmo allowlists within the config. Only accepts subcommands.",
	}

	registerSubcommands(allowlistCmd)

	return allowlistCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		deployment.GetCommand(),
		show.GetCommand(),
	)
}
