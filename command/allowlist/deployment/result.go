package deployment

import (
	"bytes"
	"fmt"

	"tynmo/types"
)

type DeploymentResult struct {
	AddAddresses    []types.Address `json:"addAddress,omitempty"`
	RemoveAddresses []types.Address `json:"removeAddress,omitempty"`
	Allowlist       []types.Address `json:"allowlist"`
}

func (r *DeploymentResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[CONTRACT DEPLOYMENT WHITELIST]\n\n")

	if len(r.AddAddresses) != 0 {
		buffer.WriteString(fmt.Sprintf("Added addresses: %s,\n", r.AddAddresses))
	}

	if len(r.RemoveAddresses) != 0 {
		buffer.WriteString(fmt.Sprintf("Removed addresses: %s,\n", r.RemoveAddresses))
	}

	buffer.WriteString(fmt.Sprintf("Contract deployment allowlist : %s,\n", r.Allowlist))

	return buffer.String()
}
