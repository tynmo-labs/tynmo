package delegatee_epoch

import (
	"bytes"
	"fmt"
	"math/big"
)

type delegateeResult struct {
	EndEpoch       *big.Int `json:"end_epoch"`
	PublicAddress  string   `json:"public_address"`
	TxHashReturned string   `json:"-"`
}

func (r *delegateeResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Set Delegatee end epoch]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *delegateeResult) Message() string {
	return fmt.Sprintf(
		"Successfully set delegatee end epoch [%v] the address [%s] with transaction hash [%s]",
		r.EndEpoch,
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *delegateeResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
