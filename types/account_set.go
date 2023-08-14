package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"tynmo/types/bitmap"
)

// AccountSet is a type alias for slice of ValidatorMetadata instances
type AccountSet []*ValidatorMetadata

func NewAccountSet() AccountSet {
	return make([]*ValidatorMetadata, 0, 0)
}

// Equals compares checks if two AccountSet instances are equal (ordering is important)
func (as AccountSet) Equals(other AccountSet) bool {
	if len(as) != len(other) {
		return false
	}

	for i := range as {
		if !as[i].Equals(other[i]) {
			return false
		}
	}

	return true
}

// fmt.Stringer implementation
func (as AccountSet) String() string {
	var buf bytes.Buffer
	for _, v := range as {
		buf.WriteString(fmt.Sprintf("%s\n", v.String()))
	}

	return buf.String()
}

// GetAddresses aggregates addresses for given AccountSet
func (as AccountSet) GetAddresses() []Address {
	res := make([]Address, 0, len(as))
	for _, account := range as {
		res = append(res, account.Address)
	}

	return res
}

// GetAddressesAsSet aggregates addresses as map for given AccountSet
func (as AccountSet) GetAddressesAsSet() map[Address]struct{} {
	res := make(map[Address]struct{}, len(as))
	for _, account := range as {
		res[account.Address] = struct{}{}
	}

	return res
}

// Len returns length of AccountSet
func (as AccountSet) Len() int {
	return len(as)
}

// ContainsNodeID checks whether ValidatorMetadata with given nodeID is present in the AccountSet
func (as AccountSet) ContainsNodeID(nodeID string) bool {
	for _, validator := range as {
		if validator.Address.String() == nodeID {
			return true
		}
	}

	return false
}

// ContainsAddress checks whether ValidatorMetadata with given address is present in the AccountSet
func (as AccountSet) ContainsAddress(address Address) bool {
	return as.Index(address) != -1
}

// Index returns index of the given ValidatorMetadata, identified by address within the AccountSet.
// If given ValidatorMetadata is not present, it returns -1.
func (as AccountSet) Index(addr Address) int {
	for indx, validator := range as {
		if validator.Address == addr {
			return indx
		}
	}

	return -1
}

// Copy returns deep copy of AccountSet
func (as AccountSet) Copy() AccountSet {
	copiedAccs := make([]*ValidatorMetadata, as.Len())
	for i, acc := range as {
		copiedAccs[i] = acc.Copy()
	}

	return AccountSet(copiedAccs)
}

// GetValidatorMetadata tries to retrieve validator account metadata by given address from the account set.
// It returns nil if such account is not found.
func (as AccountSet) GetValidatorMetadata(address Address) *ValidatorMetadata {
	i := as.Index(address)
	if i == -1 {
		return nil
	}

	return as[i]
}

// GetFilteredValidators returns filtered validators based on provided bitmap.
// Filtered validators will contain validators whose index corresponds
// to the position in bitmap which has value set to 1.
func (as AccountSet) GetFilteredValidators(bitmap bitmap.Bitmap) (AccountSet, error) {
	var filteredValidators AccountSet
	if len(as) == 0 {
		return filteredValidators, nil
	}

	if bitmap.Len() > uint64(len(as)) {
		for i := len(as); i < int(bitmap.Len()); i++ {
			if bitmap.IsSet(uint64(i)) {
				return filteredValidators, errors.New("invalid bitmap filter provided")
			}
		}
	}

	for i, validator := range as {
		if bitmap.IsSet(uint64(i)) {
			filteredValidators = append(filteredValidators, validator)
		}
	}

	return filteredValidators, nil
}

// Marshal marshals AccountSet to JSON
func (as AccountSet) Marshal() ([]byte, error) {
	return json.Marshal(as)
}

// Unmarshal unmarshals AccountSet from JSON
func (as *AccountSet) Unmarshal(b []byte) error {
	return json.Unmarshal(b, as)
}

// GetTotalVotingPower calculates sum of voting power for each validator in the AccountSet
func (as *AccountSet) GetTotalVotingPower() *big.Int {
	totalVotingPower := big.NewInt(0)
	for _, v := range *as {
		totalVotingPower = totalVotingPower.Add(totalVotingPower, v.VotingPower)
	}

	return totalVotingPower
}

func AppendAccountSet(as *AccountSet, vs *ValidatorMetadata) {
	*as = append(*as, vs)
}
