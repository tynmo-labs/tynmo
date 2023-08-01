package start

import (
	"fmt"

	"github.com/spf13/cobra"
	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/command/start/config"
	"tynmo/command/start/export"
	"tynmo/server"
	"tynmo/validators"
)

func GetCommand() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:     "start",
		Short:   "The default command that starts the Tynmo client, by bootstrapping all modules together",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	helper.RegisterGRPCAddressFlag(serverCmd)
	helper.RegisterLegacyGRPCAddressFlag(serverCmd)
	helper.RegisterJSONRPCFlag(serverCmd)

	registerSubcommands(serverCmd)
	setFlags(serverCmd)

	return serverCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// server export
		export.GetCommand(),
	)
}

func setFlags(cmd *cobra.Command) {
	defaultConfig := config.DefaultConfig()

	cmd.Flags().StringArrayVar(
		&params.premine,
		premineFlag,
		[]string{},
		fmt.Sprintf(
			"the premined accounts and balances (format: <address>:<balance>). Default premined balance: %s",
			command.DefaultPremineBalance,
		),
	)

	cmd.Flags().StringVar(
		&params.rawIBFTValidatorType,
		command.IBFTValidatorTypeFlag,
		string(validators.BLSValidatorType),
		"the type of validators in IBFT",
	)

	cmd.Flags().Uint64Var(
		&params.chainID,
		ChainIdFlag,
		command.DefaultChainID,
		"the ID of the chain",
	)

	cmd.Flags().StringVar(
		&params.initialTrieRoot,
		InitialRootFlag,
		"",
		"initial root",
	)

	cmd.Flags().StringArrayVar(
		&params.bootnodes,
		command.BootnodeFlag,
		[]string{},
		"multiAddr URL for p2p discovery bootstrap. This flag can be used multiple times",
	)

	cmd.Flags().StringArrayVar(
		&params.ibftValidatorsRaw,
		command.IBFTValidatorFlag,
		[]string{},
		"addresses to be used as IBFT validators, can be used multiple times. "+
			"Needs to be present if ibft-validators-prefix-path is omitted",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.LogLevel,
		command.LogLevelFlag,
		defaultConfig.LogLevel,
		"the log level for console output",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.GenesisPath,
		genesisPathFlag,
		defaultConfig.GenesisPath,
		"the genesis file used for starting the chain",
	)

	cmd.Flags().StringVar(
		&params.configPath,
		configFlag,
		"",
		"the path to the CLI config. Supports .json and .hcl",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.DataDir,
		dataDirFlag,
		defaultConfig.DataDir,
		"the data directory used for storing Tynmo client data",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.Libp2pAddr,
		libp2pAddressFlag,
		defaultConfig.Network.Libp2pAddr,
		"the address and port for the libp2p service",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Telemetry.PrometheusAddr,
		prometheusAddressFlag,
		"",
		"the address and port for the prometheus instrumentation service (address:port). "+
			"If only port is defined (:port) it will bind to 0.0.0.0:port",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.NatAddr,
		natFlag,
		"",
		"the external IP address without port, as can be seen by peers",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.DNSAddr,
		dnsFlag,
		"",
		"the host DNS address which can be used by a remote peer for connection",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.BlockGasTarget,
		blockGasTargetFlag,
		defaultConfig.BlockGasTarget,
		"the target block gas limit for the chain. If omitted, the value of the parent block is used",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.SecretsConfigPath,
		secretsConfigFlag,
		"",
		"the path to the SecretsManager config file. Used for Hashicorp Vault. "+
			"If omitted, the local FS secrets manager is used",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.RestoreFile,
		restoreFlag,
		"",
		"the path to the archive blockchain data to restore on initialization",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.ShouldSeal,
		sealFlag,
		defaultConfig.ShouldSeal,
		"the flag indicating that the client should seal blocks",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.Network.NoDiscover,
		command.NoDiscoverFlag,
		defaultConfig.Network.NoDiscover,
		"prevent the client from discovering other peers",
	)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxPeers,
		maxPeersFlag,
		-1,
		"the client's max number of peers allowed",
	)
	// override default usage value
	cmd.Flag(maxPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxPeers)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxInboundPeers,
		maxInboundPeersFlag,
		-1,
		"the client's max number of inbound peers allowed",
	)
	// override default usage value
	cmd.Flag(maxInboundPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxInboundPeers)
	cmd.MarkFlagsMutuallyExclusive(maxPeersFlag, maxInboundPeersFlag)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxOutboundPeers,
		maxOutboundPeersFlag,
		-1,
		"the client's max number of outbound peers allowed",
	)
	// override default usage value
	cmd.Flag(maxOutboundPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxOutboundPeers)
	cmd.MarkFlagsMutuallyExclusive(maxPeersFlag, maxOutboundPeersFlag)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.PriceLimit,
		priceLimitFlag,
		defaultConfig.TxPool.PriceLimit,
		fmt.Sprintf(
			"the minimum gas price limit to enforce for acceptance into the pool (default %d)",
			defaultConfig.TxPool.PriceLimit,
		),
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.MaxSlots,
		maxSlotsFlag,
		defaultConfig.TxPool.MaxSlots,
		"maximum slots in the pool",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.MaxAccountEnqueued,
		maxEnqueuedFlag,
		defaultConfig.TxPool.MaxAccountEnqueued,
		"maximum number of enqueued transactions per account",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.BlockTime,
		blockTimeFlag,
		defaultConfig.BlockTime,
		"minimum block time in seconds (at least 1s)",
	)

	cmd.Flags().StringArrayVar(
		&params.corsAllowedOrigins,
		corsOriginFlag,
		defaultConfig.Headers.AccessControlAllowOrigins,
		"the CORS header indicating whether any JSON-RPC response can be shared with the specified origin",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.JSONRPCBatchRequestLimit,
		jsonRPCBatchRequestLimitFlag,
		defaultConfig.JSONRPCBatchRequestLimit,
		"max length to be considered when handling json-rpc batch requests, value of 0 disables it",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.JSONRPCBlockRangeLimit,
		jsonRPCBlockRangeLimitFlag,
		defaultConfig.JSONRPCBlockRangeLimit,
		"max block range to be considered when executing json-rpc requests "+
			"that consider fromBlock/toBlock values (e.g. eth_getLogs), value of 0 disables it",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.LogFilePath,
		logFileLocationFlag,
		defaultConfig.LogFilePath,
		"write all logs to the file at specified location instead of writing them to console",
	)

	setLegacyFlags(cmd)

	setDevFlags(cmd)
}

// setLegacyFlags sets the legacy flags to preserve backwards compatibility
// with running partners
func setLegacyFlags(cmd *cobra.Command) {
	// Legacy IBFT base timeout flag
	cmd.Flags().Uint64Var(
		&params.ibftBaseTimeoutLegacy,
		ibftBaseTimeoutFlagLEGACY,
		0,
		"",
	)

	_ = cmd.Flags().MarkHidden(ibftBaseTimeoutFlagLEGACY)
}

func setDevFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&params.isDevMode,
		devFlag,
		false,
		"should the client start in dev mode (default false)",
	)

	_ = cmd.Flags().MarkHidden(devFlag)

	cmd.Flags().Uint64Var(
		&params.devInterval,
		devIntervalFlag,
		0,
		"the client's dev notification interval in seconds (default 1)",
	)

	_ = cmd.Flags().MarkHidden(devIntervalFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	// Set the grpc and json ip:port bindings
	// The config file will have precedence over --flag
	params.setRawGRPCAddress(helper.GetGRPCAddress(cmd))
	params.setRawJSONRPCAddress(helper.GetJSONRPCAddress(cmd))
	params.setJSONLogFormat(helper.GetJSONLogFormat(cmd))

	// Check if the config file has been specified
	// Config file settings will override JSON-RPC and GRPC address values
	if isConfigFileSpecified(cmd) {
		if err := params.initConfigFromFile(); err != nil {
			return err
		}
	}

	if err := params.initRawParams(); err != nil {
		return err
	}

	return nil
}

func isConfigFileSpecified(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(configFlag)
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)

	if err := runServerLoop(params.generateConfig(), outputter); err != nil {
		outputter.SetError(err)
		outputter.WriteOutput()

		return
	}
}

func runServerLoop(
	config *server.Config,
	outputter command.OutputFormatter,
) error {
	serverInstance, err := server.NewServer(config)
	if err != nil {
		return err
	}

	return helper.HandleSignals(serverInstance.Close, outputter)
}
