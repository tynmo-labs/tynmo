package unstake_delegator

import (
	"bytes"
	"fmt"
)

type StakeResult struct {
	DelegateeAddress string `json:"-"`
	PublicAddress    string `json:"-"`
	TxHashReturned   string `json:"-"`
}

func (r *StakeResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[UnStake at delegatee]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *StakeResult) Message() string {
	return fmt.Sprintf(
		"Successfully unstaked at the deletetee [%s] for address [%s] with transaction hash [%s]",
		r.DelegateeAddress,
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *StakeResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
