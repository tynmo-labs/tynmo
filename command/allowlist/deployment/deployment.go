package deployment

import (
	"fmt"

	"tynmo/command"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	deploymentCmd := &cobra.Command{
		Use:     "deployment",
		Short:   "Top level command for updating smart contract deployment allowlist. Only accepts subcommands",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(deploymentCmd)

	return deploymentCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.genesisPath,
		chainFlag,
		fmt.Sprintf("./%s", command.DefaultGenesisFileName),
		"the genesis file to update",
	)
	cmd.Flags().StringArrayVar(
		&params.addAddressRaw,
		addAddressFlag,
		[]string{},
		"adds a new address to the contract deployment allowlist",
	)

	cmd.Flags().StringArrayVar(
		&params.removeAddressRaw,
		removeAddressFlag,
		[]string{},
		"removes a new address from the contract deployment allowlist",
	)
}

func runPreRun(_ *cobra.Command, _ []string) error {
	return params.initRawParams()
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	if err := params.updateGenesisConfig(); err != nil {
		outputter.SetError(err)

		return
	}

	if err := params.overrideGenesisConfig(); err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(params.getResult())
}
