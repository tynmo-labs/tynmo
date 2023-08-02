package unstaking

import (
	"bytes"
	"fmt"
)

type IBFTUnstakeResult struct {
	PublicAddress  string `json:"-"`
	TxHashReturned string `json:"-"`
}

func (r *IBFTUnstakeResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[IBFT Unstake for Validator]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *IBFTUnstakeResult) Message() string {
	return fmt.Sprintf(
		"Successfully unstaked for the validator at address [%s] with transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *IBFTUnstakeResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
