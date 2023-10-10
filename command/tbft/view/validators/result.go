package validators

import (
	"bytes"
	"fmt"

	"tynmo/command/helper"
)

type viewResult struct {
	Validators []string `json:"validators"`
}

func (r *viewResult) Decode(infos map[string]interface{}) {

}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[VALIDATORS]\n")

	vals := make([]string, len(r.Validators))
	for i, addr := range r.Validators {
		vals[i] = fmt.Sprintf("Validators address|%s", addr)
	}
	buffer.WriteString(helper.FormatKV(vals))
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *viewResult) Message() string {
	return fmt.Sprintf(
		"Validators address [%v]",
		r.Validators,
	)
}

func (r *viewResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
