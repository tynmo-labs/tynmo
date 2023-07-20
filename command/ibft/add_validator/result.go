package add_validator

import (
	"bytes"
	"fmt"
)

type IBFTAddValidatorResult struct {
	Address string `json:"-"`
	Vote    string `json:"-"`
}

func (r *IBFTAddValidatorResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[IBFT AddValidator]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *IBFTAddValidatorResult) Message() string {
	if r.Vote == authVote {
		return fmt.Sprintf(
			"Successfully voted for the addition of address [%s] to the validator set",
			r.Address,
		)
	}

	return fmt.Sprintf(
		"Successfully voted for the removal of validator at address [%s] from the validator set",
		r.Address,
	)
}

func (r *IBFTAddValidatorResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
