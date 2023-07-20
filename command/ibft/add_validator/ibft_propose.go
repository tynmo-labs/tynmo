package add_validator

import (
	"fmt"

	"github.com/spf13/cobra"
	"tynmo/command"

	"tynmo/command/helper"
)

func GetCommand() *cobra.Command {
	ibftSnapshotCmd := &cobra.Command{
		Use:     "add_validator",
		Short:   "add a new candidate to validator set",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(ibftSnapshotCmd)

	helper.SetRequiredFlags(ibftSnapshotCmd, params.getRequiredFlags())

	return ibftSnapshotCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.addressRaw,
		addressFlag,
		"",
		"the address of the account to be voted for",
	)

	cmd.Flags().StringVar(
		&params.rawBLSPublicKey,
		blsFlag,
		"",
		"the BLS Public Key of the account to be voted for",
	)

	cmd.Flags().StringVar(
		&params.fromRaw,
		fromFlag,
		"",
		"the height to switch the new type",
	)

	cmd.Flags().StringVar(
		&params.vote,
		voteFlag,
		"",
		fmt.Sprintf(
			"requested change to the validator set. Possible values: [%s, %s]",
			authVote,
			dropVote,
		),
	)

	cmd.MarkFlagsRequiredTogether(addressFlag, voteFlag)
}

func runPreRun(_ *cobra.Command, _ []string) error {
	if err := params.validateFlags(); err != nil {
		return err
	}

	return params.initRawParams()
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	if err := params.AddValidatorCandidate(helper.GetGRPCAddress(cmd)); err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(params.getResult())
}
