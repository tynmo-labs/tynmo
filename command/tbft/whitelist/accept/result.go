package accept

import (
	"bytes"
	"fmt"
)

type acceptWhitelistResult struct {
	PublicAddress  string `json:"-"`
	TxHashReturned string `json:"-"`
}

func (r *acceptWhitelistResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Accept whitelist ownership]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *acceptWhitelistResult) Message() string {
	return fmt.Sprintf(
		"Successfully accept whitelist ownership for the address [%v] with transaction hash [%s]",
		r.PublicAddress,
		r.TxHashReturned,
	)
}

func (r *acceptWhitelistResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
