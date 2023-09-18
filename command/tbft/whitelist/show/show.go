package show

import (
	"errors"

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
var errUnexpectedResp = errors.New("unexpected response")

func GetCommand() *cobra.Command {
	whitelistCmd := &cobra.Command{
		Use:     "show",
		Short:   "show whitelist accounts",
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

	list, err := runGet(c, blockNumber)
	if err != nil {
		outputter.SetError(err)
		return
	}

	outputter.SetCommandResult(params.getResult(list))
}

func runGet(c *contract.Contract, number ethgo.BlockNumber) ([]string, error) {
	method := common.FnGetWhitelist
	res, err := c.Call(method, number)
	if err != nil {
		return nil, err
	}

	list, ok := res["0"]
	if !ok {
		return nil, errUnexpectedResp
	}

	if addresses, ok := list.([]ethgo.Address); ok {
		ret := []string{}
		for _, addr := range addresses {
			ret = append(ret, addr.String())
		}
		return ret, nil
	}
	return nil, errUnexpectedResp
}
