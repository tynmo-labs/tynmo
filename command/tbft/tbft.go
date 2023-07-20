package tbft

import (
	"github.com/spf13/cobra"
	"tynmo/command/tbft/common"
	"tynmo/command/tbft/staking"
	"tynmo/command/tbft/unstaking"
)

func GetCommand() *cobra.Command {
	tbftCmd := &cobra.Command{
		Use:   "tbft",
		Short: "Top level TynmoBFT command for interacting with the smart contract. Only accepts subcommands.",
	}

	common.RegisterJSONRPCFlag(tbftCmd)
	registerSubcommands(tbftCmd)

	return tbftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		staking.GetCommand(),
		unstaking.GetCommand(),
	)
}
