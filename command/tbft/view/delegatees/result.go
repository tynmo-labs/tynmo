package delegatees

import (
	"bytes"
	"fmt"

	"tynmo/command/helper"
)

type viewResult struct {
	Delegatees []string `json:"delagatees"`
}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[DELEGATEES]\n")

	vals := make([]string, len(r.Delegatees))
	for i, addr := range r.Delegatees {
		vals[i] = fmt.Sprintf("Delegatees address|%s", addr)
	}
	buffer.WriteString(helper.FormatKV(vals))
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *viewResult) Message() string {
	return fmt.Sprintf(
		"Current Delegatees [%v]",
		r.Delegatees,
	)
}

func (r *viewResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
