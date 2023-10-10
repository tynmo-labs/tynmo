package delegatee_info

import (
	"errors"
	"fmt"

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
		Use:     "delegatee-info",
		Short:   "show delegatee info",
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
		&params.delegatee,
		DelegateeFlag,
		"",
		DelegateeDesc,
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

	list, err := runGet(c, blockNumber, params.delegatee)
	if err != nil {
		outputter.SetError(err)
		return
	}

	outputter.SetCommandResult(params.getResult(list))
}

func runGet(c *contract.Contract, number ethgo.BlockNumber, delegatee string) (map[string]interface{}, error) {
	method := "delegateeInfo"
	res, err := c.Call(method, number, delegatee)
	if err != nil {
		return nil, err
	}

	fmt.Printf("#### %+v\n", res)

	resp, ok := res["0"]
	if !ok {
		return nil, errUnexpectedResp
	}

	if info, ok := resp.(map[string]interface{}); ok {
		return info, nil
	}

	return nil, errUnexpectedResp
}
