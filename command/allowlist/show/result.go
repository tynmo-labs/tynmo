package show

import (
	"bytes"
	"fmt"
)

type ShowResult struct {
	Allowlists Allowlists
}

func (r *ShowResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[WHITELISTS]\n\n")

	buffer.WriteString(fmt.Sprintf("Contract deployment allowlist : %s,\n", r.Allowlists.deployment))

	return buffer.String()
}
