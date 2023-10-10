package delegatee_info

import (
	"bytes"
	"fmt"

	"github.com/umbracle/ethgo"
)

type viewResult struct {
	Info       string   `json:"info"`
	Delegators []string `json:"delagators"`
	EndEpoch   uint64   `json:"endEpoch"`
	Percetage  uint64   `json:"percentage"`
	Reaward    uint64   `json:"reward"`
	Stake      uint64   `json:"stake"`
}

func (r *viewResult) Decode(infos map[string]interface{}) {
	if delegators, has := infos["delagators"]; has && delegators != nil {
		if addresses, ok := delegators.([]ethgo.Address); ok {
			ret := []string{}
			for _, addr := range addresses {
				ret = append(ret, addr.String())
			}
			r.Delegators = ret
		}
	}
}

func (r *viewResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[DELEGATEE INFO]\n")
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
