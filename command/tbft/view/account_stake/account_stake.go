package account_stake

import (
	"errors"
	"fmt"
	"math/big"

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

var params viewParams
var errUnexpectedResp = errors.New("unexpected response")

func GetCommand() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:     "account-stake",
		Short:   "show account stake",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(viewCmd)

	return viewCmd
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
		&params.account,
		AccountFlag,
		"",
		AccountDesc,
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

	blockNumber, err := common.BlockNumber(client)
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

	amount, err := runGet(c, blockNumber, params.account)
	if err != nil {
		outputter.SetError(err)
		return
	}

	outputter.SetCommandResult(params.getResult(amount))
}

func runGet(c *contract.Contract, number ethgo.BlockNumber, account string) (*big.Int, error) {
	method := "accountStake"
	res, err := c.Call(method, number, account)
	if err != nil {
		return nil, err
	}

	fmt.Printf("#### %+v\n", res)

	resp, ok := res["0"]
	if !ok {
		return nil, errUnexpectedResp
	}

	if amount, ok := resp.(*big.Int); ok {
		return amount, nil
	}

	return nil, errUnexpectedResp
}
