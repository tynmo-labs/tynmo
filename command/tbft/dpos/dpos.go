package dpos

import (
	"github.com/spf13/cobra"

	"tynmo/command/tbft/common"
	"tynmo/command/tbft/dpos/delegatee_epoch"
	"tynmo/command/tbft/dpos/delegatee_percentage"
	"tynmo/command/tbft/dpos/register_delegatee"
	"tynmo/command/tbft/dpos/stake_delegator"
	"tynmo/command/tbft/dpos/unstake_delegator"
	"tynmo/command/tbft/dpos/withdraw_delegator_stake"
)

func GetCommand() *cobra.Command {
	tbftCmd := &cobra.Command{
		Use:   "dpos",
		Short: "Subcommand for interacting with dpos contract.",
	}

	common.RegisterJSONRPCFlag(tbftCmd)
	registerSubcommands(tbftCmd)

	return tbftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		register_delegatee.GetCommand(),
		stake_delegator.GetCommand(),
		unstake_delegator.GetCommand(),
		withdraw_delegator_stake.GetCommand(),
		delegatee_epoch.GetCommand(),
		delegatee_percentage.GetCommand(),
	)
}
