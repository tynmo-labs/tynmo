package common

import (
	"fmt"

	"github.com/spf13/cobra"
	"tynmo/command"
	"tynmo/command/helper"
)

const (
	DefaultJSONRPCPort int = 10002
)

// RegisterJSONRPCFlag registers the base JSON-RPC address flag for all child commands
func RegisterJSONRPCFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		command.JSONRPCFlag,
		fmt.Sprintf("http://%s:%d", helper.LocalHostBinding, DefaultJSONRPCPort),
		"the JSON-RPC interface",
	)
}
