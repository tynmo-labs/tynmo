package unstaking

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

var params unstakeParams

func GetCommand() *cobra.Command {
	unstakeCmd := &cobra.Command{
		Use:     "unstake",
		Short:   "Unstakes the amount sent for a validator",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(unstakeCmd)

	return unstakeCmd
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
		&params.amount,
		AmountFlag,
		"",
		AmountToUnstakeDesc,
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

	txn, err := c.Txn("withdrawStake", params.amountValue)
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
