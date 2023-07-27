package validators

import (
	"fmt"
	"math/big"
	"tynmo/types"
	"github.com/umbracle/fastrlp"
)

// ECDSAValidator is a validator using ECDSA signing algorithm
type ECDSAValidator struct {
	Address types.Address
	Stake   *big.Int
	Balance *big.Int
}

// NewECDSAValidator is a constructor of ECDSAValidator
func NewECDSAValidator(addr types.Address, args ...interface{}) *ECDSAValidator {
	var accountStake = big.NewInt(0)
	var accountBalance = big.NewInt(0)
	argsLen := len(args)

	if argsLen >= 1 {
		ok := false
		accountStake, ok = args[0].(*big.Int)
		if !ok {
			return nil
		}
	}
	if argsLen == 2 {
		ok := false
		accountBalance, ok = args[1].(*big.Int)
		if !ok {
			return nil
		}
	}
	if argsLen > 2 {
		return nil
	}

	return &ECDSAValidator{
		Address: addr,
		Stake:   accountStake,
		Balance: accountBalance,
	}
}

// Type returns the ValidatorType of ECDSAValidator
func (v *ECDSAValidator) Type() ValidatorType {
	return ECDSAValidatorType
}

// String returns string representation of ECDSAValidator
func (v *ECDSAValidator) String() string {
	return v.Address.String()
}

// Addr returns the validator address
func (v *ECDSAValidator) Addr() types.Address {
	return v.Address
}

// Copy returns copy of ECDSAValidator
func (v *ECDSAValidator) Copy() Validator {
	return &ECDSAValidator{
		Address: v.Address,
		Balance: v.Balance,
		Stake:   v.Stake,
	}
}

// Equal checks the given validator matches with its data
func (v *ECDSAValidator) Equal(vr Validator) bool {
	vv, ok := vr.(*ECDSAValidator)
	if !ok {
		return false
	}

	return v.Address == vv.Address
}

// MarshalRLPWith is a RLP Marshaller
func (v *ECDSAValidator) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewBytes(v.Address.Bytes()))
	vv.Set(arena.NewBigInt(v.Stake))
	vv.Set(arena.NewBigInt(v.Balance))

	return vv
}

// UnmarshalRLPFrom is a RLP Unmarshaller
func (v *ECDSAValidator) UnmarshalRLPFrom(p *fastrlp.Parser, val *fastrlp.Value) error {
	elems, err := val.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 2 {
		return fmt.Errorf("incorrect number of elements to decode BLSValidator, expected 3 but found %d", len(elems))
	}

	if err := elems[0].GetAddr(v.Address[:]); err != nil {
		return fmt.Errorf("failed to decode Address: %w", err)
	}

	if v.Stake == nil {
		v.Stake = new(big.Int)
	}
	if err = elems[1].GetBigInt(v.Stake); err != nil {
		return fmt.Errorf("failed to decode Stake: %w", err)
	}

	if v.Balance == nil {
		v.Balance = new(big.Int)
	}
	if err = elems[2].GetBigInt(v.Balance); err != nil {
		return fmt.Errorf("failed to decode Balance: %w", err)
	}

	return nil
}

// Bytes returns bytes of ECDSAValidator
func (v *ECDSAValidator) Bytes() []byte {
	return v.Address.Bytes()
}

// SetFromBytes parses given bytes
func (v *ECDSAValidator) SetFromBytes(input []byte) error {
	return types.UnmarshalRlp(v.UnmarshalRLPFrom, input)
}

// GetBalance returns the validator balance
func (v *ECDSAValidator) GetBalance() *big.Int {
	return v.Balance
}

// GetStake returns the validator stake
func (v *ECDSAValidator) GetStake() *big.Int {
	return v.Stake
}
