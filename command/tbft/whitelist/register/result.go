package register

import (
	"bytes"
	"fmt"
)

type registerWhitelistResult struct {
	PublicAddress     string `json:"-"`
	TxHashReturned    string `json:"-"`
	KeyTxHashReturned string `json:"-"`
}

func (r *registerWhitelistResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Whitelist Register for Validator]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *registerWhitelistResult) Message() string {
	return fmt.Sprintf(
		"Successfully register for the validator at address [%s] with register transaction hash [%s] and register public key transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
		r.KeyTxHashReturned,
	)
}

func (r *registerWhitelistResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
