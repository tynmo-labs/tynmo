package ibft

import (
	"github.com/spf13/cobra"
	"tynmo/command/helper"
	"tynmo/command/ibft/add_validator"
	"tynmo/command/ibft/candidates"
	"tynmo/command/ibft/propose"
	"tynmo/command/ibft/quorum"
	"tynmo/command/ibft/snapshot"
	"tynmo/command/ibft/status"
	_switch "tynmo/command/ibft/switch"
)

func GetCommand() *cobra.Command {
	ibftCmd := &cobra.Command{
		Use:   "ibft",
		Short: "Top level IBFT command for interacting with the IBFT consensus. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(ibftCmd)

	registerSubcommands(ibftCmd)

	return ibftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// ibft status
		status.GetCommand(),
		// ibft snapshot
		snapshot.GetCommand(),
		// ibft propose
		propose.GetCommand(),
		// ibft candidates
		candidates.GetCommand(),
		// ibft switch
		_switch.GetCommand(),
		// ibft quorum
		quorum.GetCommand(),
		add_validator.GetCommand(),
	)
}
