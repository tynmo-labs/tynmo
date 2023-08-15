package validators

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"tynmo/helper/hex"
	"tynmo/types"

	"github.com/umbracle/fastrlp"
)

var (
	ErrInvalidTypeAssert = errors.New("invalid type assert")
)

type BLSValidatorPublicKey []byte

// String returns a public key in hex
func (k BLSValidatorPublicKey) String() string {
	return hex.EncodeToHex(k[:])
}

// MarshalText implements encoding.TextMarshaler
func (k BLSValidatorPublicKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText parses an BLS Public Key in hex
func (k *BLSValidatorPublicKey) UnmarshalText(input []byte) error {
	kk, err := hex.DecodeHex(string(input))
	if err != nil {
		return err
	}

	*k = kk

	return nil
}

// BLSValidator is a validator using BLS signing algorithm
type BLSValidator struct {
	Address      types.Address
	BLSPublicKey BLSValidatorPublicKey
	Stake        *big.Int
	Balance      *big.Int
}

// NewBLSValidator is a constructor of BLSValidator
func NewBLSValidator(addr types.Address, blsPubKey []byte, args ...interface{}) *BLSValidator {
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

	return &BLSValidator{
		Address:      addr,
		BLSPublicKey: blsPubKey,
		Stake:        accountStake,
		Balance:      accountBalance,
	}
}

// Type returns the ValidatorType of BLSValidator
func (v *BLSValidator) Type() ValidatorType {
	return BLSValidatorType
}

// String returns string representation of BLSValidator
// Format => [Address]:[BLSPublicKey]
func (v *BLSValidator) String() string {
	return fmt.Sprintf(
		"%s:%s",
		v.Address.String(),
		hex.EncodeToHex(v.BLSPublicKey),
	)
}

// Addr returns the validator address
func (v *BLSValidator) Addr() types.Address {
	return v.Address
}

// Copy returns copy of BLS Validator
func (v *BLSValidator) Copy() Validator {
	pubkey := make([]byte, len(v.BLSPublicKey))
	copy(pubkey, v.BLSPublicKey)

	return &BLSValidator{
		Address:      v.Address,
		BLSPublicKey: pubkey,
		Stake:        big.NewInt(0).Set(v.Stake),
		Balance:      big.NewInt(0).Set(v.Balance),
	}
}

// Equal checks the given validator matches with its data
func (v *BLSValidator) Equal(vr Validator) bool {
	vv, ok := vr.(*BLSValidator)
	if !ok {
		return false
	}

	return v.Address == vv.Address && bytes.Equal(v.BLSPublicKey, vv.BLSPublicKey)
}

// MarshalRLPWith is a RLP Marshaller
func (v *BLSValidator) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewBytes(v.Address.Bytes()))
	vv.Set(arena.NewCopyBytes(v.BLSPublicKey))
	vv.Set(arena.NewBigInt(v.Stake))
	vv.Set(arena.NewBigInt(v.Balance))

	return vv
}

// UnmarshalRLPFrom is a RLP Unmarshaller
func (v *BLSValidator) UnmarshalRLPFrom(p *fastrlp.Parser, val *fastrlp.Value) error {
	elems, err := val.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 2 {
		return fmt.Errorf("incorrect number of elements to decode BLSValidator, expected 2 but found %d", len(elems))
	}

	if err := elems[0].GetAddr(v.Address[:]); err != nil {
		return fmt.Errorf("failed to decode Address: %w", err)
	}

	if v.BLSPublicKey, err = elems[1].GetBytes(v.BLSPublicKey); err != nil {
		return fmt.Errorf("failed to decode BLSPublicKey: %w", err)
	}

	if v.Stake == nil {
		v.Stake = new(big.Int)
	}
	if err = elems[2].GetBigInt(v.Stake); err != nil {
		return fmt.Errorf("failed to decode Stake: %w", err)
	}

	if v.Balance == nil {
		v.Balance = new(big.Int)
	}
	if err = elems[3].GetBigInt(v.Balance); err != nil {
		return fmt.Errorf("failed to decode Balance: %w", err)
	}

	return nil
}

// Bytes returns bytes of BLSValidator in RLP encode
func (v *BLSValidator) Bytes() []byte {
	return types.MarshalRLPTo(v.MarshalRLPWith, nil)
}

// SetFromBytes parses given bytes in RLP encode and map to its fields
func (v *BLSValidator) SetFromBytes(input []byte) error {
	return types.UnmarshalRlp(v.UnmarshalRLPFrom, input)
}

// GetBalance returns the validator balance
func (v *BLSValidator) GetBalance() *big.Int {
	return v.Balance
}

// GetStake returns the validator stake
func (v *BLSValidator) GetStake() *big.Int {
	return v.Stake
}
