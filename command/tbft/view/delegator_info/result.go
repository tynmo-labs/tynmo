package delegator_info

import (
	"bytes"
	"fmt"
)

type viewResult struct {
	Info string `json:"info"`
}

func (r *viewResult) Decode(infos map[string]interface{}) {

}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[DELEGATOR INFO]\n")
	buffer.WriteString(r.Info)
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *viewResult) Message() string {
	return fmt.Sprintf(
		"Current Delegatee info [%v]",
		r.Info,
	)
}

func (r *viewResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
