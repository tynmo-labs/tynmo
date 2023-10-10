package delegatee_percentage

import (
	"bytes"
	"fmt"
	"math/big"
)

type delegateeResult struct {
	Percentage     *big.Int `json:"percentage"`
	PublicAddress  string   `json:"-"`
	TxHashReturned string   `json:"-"`
}

func (r *delegateeResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Set Delegatee end percentage]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *delegateeResult) Message() string {
	return fmt.Sprintf(
		"Successfully Set Delegatee end percentage [%v] for the address [%s] with transaction hash [%s]",
		r.Percentage,
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *delegateeResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
