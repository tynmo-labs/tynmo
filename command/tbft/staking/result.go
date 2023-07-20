package staking

import (
	"bytes"
	"fmt"
)

type IBFTStakeResult struct {
	PublicAddress  string `json:"-"`
	TxHashReturned string `json:"-"`
}

func (r *IBFTStakeResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[IBFT Stake for Validator]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *IBFTStakeResult) Message() string {
	return fmt.Sprintf(
		"Successfully staked for the validator at address [%s] with transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *IBFTStakeResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
