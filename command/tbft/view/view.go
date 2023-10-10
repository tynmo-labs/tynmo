package view

import (
	"github.com/spf13/cobra"

	"tynmo/command/tbft/common"
	"tynmo/command/tbft/view/account_stake"
	"tynmo/command/tbft/view/delegatee_info"
	"tynmo/command/tbft/view/delegatees"
	"tynmo/command/tbft/view/delegator_info"
	"tynmo/command/tbft/view/delegators_in_delegatee"
	"tynmo/command/tbft/view/validators"
)

func GetCommand() *cobra.Command {
	tbftCmd := &cobra.Command{
		Use:   "view",
		Short: "Subcommand for interacting with tbft contract.",
	}

	common.RegisterJSONRPCFlag(tbftCmd)
	registerSubcommands(tbftCmd)

	return tbftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		delegatees.GetCommand(),
		delegators_in_delegatee.GetCommand(),
		delegatee_info.GetCommand(),
		delegator_info.GetCommand(),
		account_stake.GetCommand(),
		validators.GetCommand(),
	)
}
