package contractsapi

import (
	"math/big"

	"tynmo/contracts/abis"
	"tynmo/types"

	"github.com/umbracle/ethgo/abi"
)

const (
	//methodValidators             = "validators"
	//methodValidatorBLSPublicKeys = "validatorBLSPublicKeys"
	methodAccountStake = "accountStake"
)

type SCAccountStakeFn struct {
	Addr types.Address `abi:"addr"`
}

func (c *SCAccountStakeFn) Sig() []byte {
	return abis.StakingABI.Methods[methodAccountStake].ID()
}

func (c *SCAccountStakeFn) EncodeAbi() ([]byte, error) {
	return abis.StakingABI.Methods[methodAccountStake].Encode(c)
}

func (c *SCAccountStakeFn) DecodeAbi(buf []byte) error {
	return decodeMethod(abis.StakingABI.Methods[methodAccountStake], buf, c)
}

type Validator struct {
	Address     types.Address `abi:"_address"`
	BlsKey      [4]*big.Int   `abi:"blsKey"`
	VotingPower *big.Int      `abi:"votingPower"`
}

var ValidatorABIType = abi.MustNewType("tuple(address _address,uint256[4] blsKey,uint256 votingPower)")

func (v *Validator) EncodeAbi() ([]byte, error) {
	return ValidatorABIType.Encode(v)
}

func (v *Validator) DecodeAbi(buf []byte) error {
	return decodeStruct(ValidatorABIType, buf, &v)
}
