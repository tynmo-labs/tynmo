package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

var errUnexpectedResp = errors.New("unexpected response")

const (
	FnRegisterBLSPublicKey = "registerBLSPublicKey"
	FnAddWhitelist         = "whitelistValidators"
	FnGetWhitelist         = "whitelist"
	FnRegister             = "registerValidator"
	FnTransferOwnership    = "transferOwnership"
	FnAcceptOwnership      = "acceptOwnership"
)

func BlockNumber(client *jsonrpc.Client) (ethgo.BlockNumber, error) {
	var out string
	if err := client.Call("eth_blockNumber", &out); err != nil {
		return 0, err
	}
	number, err := parseUint64orHex(out)
	if err != nil {
		return 0, err
	}
	return ethgo.BlockNumber(number), nil
}

func parseUint64orHex(str string) (uint64, error) {
	base := 10
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
		base = 16
	}
	return strconv.ParseUint(str, base, 64)
}

func DecodeCallResponse(resp map[string]interface{}) (interface{}, error) {
	msg, ok := resp["0"]
	if ok {
		return msg, nil
	}
	return nil, errUnexpectedResp
}
