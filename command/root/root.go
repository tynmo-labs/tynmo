package root

import (
	"fmt"
	"os"

	"tynmo/command/backup"
	"tynmo/command/genesis"
	"tynmo/command/helper"
	"tynmo/command/ibft"
	"tynmo/command/license"
	"tynmo/command/monitor"
	"tynmo/command/peers"
	"tynmo/command/secrets"
	"tynmo/command/server"
	"tynmo/command/start"
	"tynmo/command/status"
	"tynmo/command/tbft"
	"tynmo/command/txpool"
	"tynmo/command/version"
	"tynmo/command/whitelist"

	"github.com/spf13/cobra"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Tynmo is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		start.GetCommand(),
		whitelist.GetCommand(),
		license.GetCommand(),
		tbft.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
