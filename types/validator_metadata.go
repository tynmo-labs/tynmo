package types

import (
	"fmt"
	"math/big"
)

// ValidatorMetadata represents a validator metadata (its public identity)
type ValidatorMetadata struct {
	Address     Address
	VotingPower *big.Int
	IsActive    bool
	IsSealer    bool
}

func NewValidatorMetadata(addr Address, vp *big.Int, isActive, isSealer bool) *ValidatorMetadata {
	return &ValidatorMetadata{
		Address:     addr,
		VotingPower: vp,
		IsActive:    isActive,
		IsSealer:    isSealer,
	}
}

// Equals checks ValidatorMetadata equality
func (v *ValidatorMetadata) Equals(b *ValidatorMetadata) bool {
	if b == nil {
		return false
	}

	return v.EqualAddress(b) && v.VotingPower.Cmp(b.VotingPower) == 0 && v.IsActive == b.IsActive
}

// EqualAddress checks ValidatorMetadata equality against Address and BlsKey fields
func (v *ValidatorMetadata) EqualAddress(b *ValidatorMetadata) bool {
	if b == nil {
		return false
	}

	return v.Address == b.Address
}

// Copy returns a deep copy of ValidatorMetadata
func (v *ValidatorMetadata) Copy() *ValidatorMetadata {
	return &ValidatorMetadata{
		Address:     BytesToAddress(v.Address[:]),
		VotingPower: new(big.Int).Set(v.VotingPower),
		IsActive:    v.IsActive,
		IsSealer:    v.IsSealer,
	}
}

// fmt.Stringer implementation
func (v *ValidatorMetadata) String() string {
	return fmt.Sprintf("Address=%v; Is Active=%v; Voting Power=%d;Is Sealer=%v;",
		v.Address.String(), v.IsActive, v.VotingPower, v.IsSealer)
}
