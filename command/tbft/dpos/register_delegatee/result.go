package register_delegatee

import (
	"bytes"
	"fmt"
)

type registerResult struct {
	PublicAddress  string `json:"-"`
	TxHashReturned string `json:"-"`
}

func (r *registerResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[register delegatee]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *registerResult) Message() string {
	return fmt.Sprintf(
		"Successfully register delegatee for the address [%s] with transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *registerResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
