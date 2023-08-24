package config

import (
	"tynmo/chain"
	"tynmo/types"
)

// GetAllowlist fetches allowlist object from the config
func GetAllowlist(config *chain.Chain) *chain.Allowlists {
	return config.Params.Allowlists
}

// GetDeploymentAllowlist fetches deployment allowlist from the genesis config
// if doesn't exist returns empty list
func GetDeploymentAllowlist(genesisConfig *chain.Chain) ([]types.Address, error) {
	// Fetch allowlist config if exists, if not init
	allowlistConfig := GetAllowlist(genesisConfig)

	// Extract deployment allowlist if exists, if not init
	if allowlistConfig == nil {
		return make([]types.Address, 0), nil
	}

	return allowlistConfig.Deployment, nil
}
