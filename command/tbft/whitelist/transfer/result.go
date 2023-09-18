package transfer

import (
	"bytes"
	"fmt"
)

type transferWhitelistResult struct {
	PublicAddress  string `json:"-"`
	TxHashReturned string `json:"-"`
}

func (r *transferWhitelistResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Transfer whitelist owner]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *transferWhitelistResult) Message() string {
	return fmt.Sprintf(
		"Successfully transfer whitelist owner to the address [%v] with transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *transferWhitelistResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
