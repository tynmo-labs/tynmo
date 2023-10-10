package delegators_in_delegatee

import (
	"bytes"
	"fmt"

	"tynmo/command/helper"
)

type viewResult struct {
	Delegators []string `json:"delegators"`
}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[Delegators]\n")

	vals := make([]string, len(r.Delegators))
	for i, addr := range r.Delegators {
		vals[i] = fmt.Sprintf("Delegators address|%s", addr)
	}
	buffer.WriteString(helper.FormatKV(vals))
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *viewResult) Message() string {
	return fmt.Sprintf(
		"Current Delegators [%v]",
		r.Delegators,
	)
}

func (r *viewResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
