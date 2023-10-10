package register_delegatee

import (
	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/contracts/abis"
	"tynmo/contracts/staking"

	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

var params registerParams

func GetCommand() *cobra.Command {
	dposCmd := &cobra.Command{
		Use:     "register-delegatee",
		Short:   "register delegatee",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(dposCmd)

	return dposCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		AccountDirFlag,
		"",
		AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.privateKeyStr,
		PrivateKeyFlag,
		"",
		PrivateKeyDesc,
	)

	cmd.Flags().StringVar(
		&params.percentage,
		PercentageFlag,
		"0",
		PercentageDesc,
	)

	cmd.Flags().StringVar(
		&params.endEpoch,
		EndEpochFlag,
		"0",
		EndEpochDesc,
	)

	cmd.MarkFlagsMutuallyExclusive(AccountDirFlag, PrivateKeyFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	if err := params.validateFlags(); err != nil {
		return err
	}

	return params.initRawParams()
}

var Number = ethgo.BlockNumber(41855)

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	abiContract := abi.MustNewABI(abis.StakingJSONABI)

	addr := ethgo.Address(staking.AddrStakingContract)

	client, err := jsonrpc.NewClient(params.jsonRPC)
	if err != nil {
		outputter.SetError(err)
		return
	}

	key := wallet.NewKey(params.privateKey)

	opts := []contract.ContractOption{
		contract.WithJsonRPC(client.Eth()),
		contract.WithSender(key),
	}
	c := contract.NewContract(addr, abiContract, opts...)

	var txn contract.Txn
	txn, err = c.Txn("registerDelegatee", params.percentageValue, params.endEpochValue)
	if err != nil {
		outputter.SetError(err)
		return
	}

	err = txn.Do()
	if err != nil {
		outputter.SetError(err)
		return
	}

	receipt, err := txn.Wait()
	if err != nil {
		outputter.SetError(err)
		return
	}

	outputter.SetCommandResult(params.getResult(receipt.TransactionHash.String()))

}
