package staking

import (
	"errors"
	"math/big"

	"github.com/umbracle/ethgo"

	"github.com/umbracle/ethgo/abi"
	"tynmo/consensus/tynmobft/contractsapi"
	"tynmo/contracts/abis"
	"tynmo/state/runtime"
	"tynmo/types"
)

const (
	methodValidators             = "validators"
	methodValidatorBLSPublicKeys = "validatorBLSPublicKeys"
	methodAccountStake           = "accountStake"
)

var (
	// staking contract address
	AddrStakingContract = types.StringToAddress("1001")

	// Gas limit used when querying the validator set
	queryGasLimit uint64 = 1000000

	ErrMethodNotFoundInABI = errors.New("method not found in ABI")
	ErrFailedTypeAssertion = errors.New("failed type assertion")
)

// TxQueryHandler is a interface to call view method in the contract
type TxQueryHandler interface {
	Apply(*types.Transaction) (*runtime.ExecutionResult, error)
	GetNonce(types.Address) uint64
}

// decodeWeb3ArrayOfBytes is a helper function to parse the data
// representing array of bytes in contract result
func decodeWeb3ArrayOfBytes(
	result interface{},
) ([][]byte, error) {
	mapResult, ok := result.(map[string]interface{})
	if !ok {
		return nil, ErrFailedTypeAssertion
	}

	bytesArray, ok := mapResult["0"].([][]byte)
	if !ok {
		return nil, ErrFailedTypeAssertion
	}

	return bytesArray, nil
}

// createCallViewTx is a helper function to create a transaction to call view method
func createCallViewTx(
	from types.Address,
	contractAddress types.Address,
	methodID []byte,
	nonce uint64,
) *types.Transaction {
	return &types.Transaction{
		From:     from,
		To:       &contractAddress,
		Input:    methodID,
		Nonce:    nonce,
		Gas:      queryGasLimit,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(0),
	}
}

// DecodeValidators parses contract call result and returns array of address
func DecodeValidators(method *abi.Method, returnValue []byte) ([]types.Address, error) {
	decodedResults, err := method.Outputs.Decode(returnValue)
	if err != nil {
		return nil, err
	}

	results, ok := decodedResults.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed type assertion from decodedResults to map")
	}

	web3Addresses, ok := results["0"].([]ethgo.Address)

	if !ok {
		return nil, errors.New("failed type assertion from results[0] to []ethgo.Address")
	}

	addresses := make([]types.Address, len(web3Addresses))
	for idx, waddr := range web3Addresses {
		addresses[idx] = types.Address(waddr)
	}

	return addresses, nil
}

// QueryValidators is a helper function to get validator addresses from contract
func QueryValidators(t TxQueryHandler, from types.Address) ([]types.Address, error) {
	method, ok := abis.StakingABI.Methods[methodValidators]
	if !ok {
		return nil, ErrMethodNotFoundInABI
	}

	res, err := t.Apply(createCallViewTx(
		from,
		AddrStakingContract,
		method.ID(),
		t.GetNonce(from),
	))

	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, res.Err
	}

	return DecodeValidators(method, res.ReturnValue)
}

// decodeBLSPublicKeys parses contract call result and returns array of bytes
func decodeBLSPublicKeys(
	method *abi.Method,
	returnValue []byte,
) ([][]byte, error) {
	decodedResults, err := method.Outputs.Decode(returnValue)
	if err != nil {
		return nil, err
	}

	blsPublicKeys, err := decodeWeb3ArrayOfBytes(decodedResults)
	if err != nil {
		return nil, err
	}

	return blsPublicKeys, nil
}

// QueryBLSPublicKeys is a helper function to get BLS Public Keys from contract
func QueryBLSPublicKeys(t TxQueryHandler, from types.Address) ([][]byte, error) {
	method, ok := abis.StakingABI.Methods[methodValidatorBLSPublicKeys]
	if !ok {
		return nil, ErrMethodNotFoundInABI
	}

	res, err := t.Apply(createCallViewTx(
		from,
		AddrStakingContract,
		method.ID(),
		t.GetNonce(from),
	))

	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, res.Err
	}

	return decodeBLSPublicKeys(method, res.ReturnValue)
}

// decodeAccountStake parses contract call result and returns uint256
func decodeAccountStake(
	method *abi.Method,
	returnValue []byte,
) (*big.Int, error) {
	decodedResults, err := method.Outputs.Decode(returnValue)
	if err != nil {
		return nil, err
	}

	results, ok := decodedResults.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed type assertion from decodedResults to map")
	}

	result, ok := results["0"].(*big.Int)

	if !ok {
		return nil, errors.New("failed type assertion from results[0] to *big.Int")
	}

	return result, nil
}

// QueryAccountStake is a helper function to get validator's stake from contract
func QueryAccountStake(t TxQueryHandler, from, accountStakeAddress types.Address) (*big.Int, error) {
	method, ok := abis.StakingABI.Methods[methodAccountStake]
	//method.Encode(accountStakeAddress)
	if !ok {
		return nil, ErrMethodNotFoundInABI
	}

	scAccountStakeFn := &contractsapi.SCAccountStakeFn{Addr: accountStakeAddress}
	scAccountStakeEncoded, err := scAccountStakeFn.EncodeAbi()
	if err != nil {
		return nil, err
	}

	res, err := t.Apply(createCallViewTx(
		from,
		AddrStakingContract,
		scAccountStakeEncoded,
		t.GetNonce(from),
	))

	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, res.Err
	}

	return decodeAccountStake(method, res.ReturnValue)
}
