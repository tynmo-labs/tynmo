package staking

import (
	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/contracts/abis"
	"tynmo/contracts/staking"
)

var params stakeParams

func GetCommand() *cobra.Command {
	stakeCmd := &cobra.Command{
		Use:     "stake",
		Short:   "Stakes the amount sent for a validator",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(stakeCmd)

	return stakeCmd
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
		AmountToStakeDesc,
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

	txn, err := c.Txn("stake")
	txn.WithOpts(&contract.TxnOpts{Value: params.amountValue})
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
