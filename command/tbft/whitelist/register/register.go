package register

import (
	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/command/tbft/whitelist/common"
	"tynmo/contracts/abis"
	"tynmo/contracts/staking"

	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

var params whitelistParams

func GetCommand() *cobra.Command {
	whitelistCmd := &cobra.Command{
		Use:     "register",
		Short:   "register for validator",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(whitelistCmd)

	return whitelistCmd
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
	txn, err = c.Txn(common.FnRegister)
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

	var txnKey contract.Txn
	txnKey, err = c.Txn(common.FnRegisterBLSPublicKey, params.blsPubKey)
	if err != nil {
		outputter.SetError(err)
		return
	}

	err = txnKey.Do()
	if err != nil {
		outputter.SetError(err)
		return
	}

	receiptKey, err := txnKey.Wait()
	if err != nil {
		outputter.SetError(err)
		return
	}

	outputter.SetCommandResult(params.getResult(receipt.TransactionHash.String(), receiptKey.TransactionHash.String()))
}
