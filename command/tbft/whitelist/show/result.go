package show

import (
	"bytes"
	"fmt"

	"tynmo/command/helper"
)

type showWhitelistResult struct {
	Validators []string `json:"Validators"`
}

func (r *showWhitelistResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[WHITELIST VALIDATORS]\n")

	vals := make([]string, len(r.Validators))
	for i, addr := range r.Validators {
		vals[i] = fmt.Sprintf("Validator address|%s", addr)
	}
	buffer.WriteString(helper.FormatKV(vals))
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *showWhitelistResult) Message() string {
	return fmt.Sprintf(
		"Current validators [%v]",
		r.Validators,
	)
}

func (r *showWhitelistResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
