package whitelist

import (
	"github.com/spf13/cobra"

	"tynmo/command/tbft/common"
	"tynmo/command/tbft/whitelist/accept"
	"tynmo/command/tbft/whitelist/add"
	"tynmo/command/tbft/whitelist/register"
	"tynmo/command/tbft/whitelist/transfer"
)

func GetCommand() *cobra.Command {
	tbftCmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Subcommand for interacting with whitelist contract.",
	}

	common.RegisterJSONRPCFlag(tbftCmd)
	registerSubcommands(tbftCmd)

	return tbftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		add.GetCommand(),
		register.GetCommand(),
		transfer.GetCommand(),
		accept.GetCommand(),
	)
}
