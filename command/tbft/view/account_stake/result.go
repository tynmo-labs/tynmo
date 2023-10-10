package account_stake

import (
	"bytes"
	"fmt"
	"math/big"
)

type viewResult struct {
	Amount *big.Int `json:"amount"`
}

func (r *viewResult) Decode(infos map[string]interface{}) {

}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Account stake amount]\n")
	buffer.WriteString(r.Amount.String())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *viewResult) Message() string {
	return fmt.Sprintf(
		"Account stake amount [%v]",
		r.Amount.String(),
	)
}

func (r *viewResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
