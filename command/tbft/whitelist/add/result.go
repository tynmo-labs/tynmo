package add

import (
	"bytes"
	"fmt"
)

type addWhitelistResult struct {
	PublicAddresses []string `json:"-"`
	TxHashReturned  string   `json:"-"`
}

func (r *addWhitelistResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Add accounts to whitelist]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *addWhitelistResult) Message() string {
	return fmt.Sprintf(
		"Successfully add whitelist for the address [%v] with transaction hash [%s]",
		r.PublicAddresses,
		r.TxHashReturned,
	)
}

func (r *addWhitelistResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
