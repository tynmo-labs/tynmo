package deployment

import (
	"fmt"
	"os"

	"tynmo/chain"
	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/helper/config"
	"tynmo/types"
)

const (
	chainFlag         = "chain"
	addAddressFlag    = "add"
	removeAddressFlag = "remove"
)

var (
	params = &deploymentParams{}
)

type deploymentParams struct {
	// raw addresses, entered by CLI commands
	addAddressRaw    []string
	removeAddressRaw []string

	// addresses, converted from raw addresses
	addAddresses    []types.Address
	removeAddresses []types.Address

	// genesis file
	genesisPath   string
	genesisConfig *chain.Chain

	// deployment allowlist from genesis configuration
	allowlist []types.Address
}

func (p *deploymentParams) initRawParams() error {
	// convert raw addresses to appropriate format
	if err := p.initRawAddresses(); err != nil {
		return err
	}

	// init genesis configuration
	if err := p.initChain(); err != nil {
		return err
	}

	return nil
}

func (p *deploymentParams) initRawAddresses() error {
	// convert addresses to be added from string to type.Address
	p.addAddresses = unmarshallRawAddresses(p.addAddressRaw)

	// convert addresses to be removed from string to type.Address
	p.removeAddresses = unmarshallRawAddresses(p.removeAddressRaw)

	return nil
}

func (p *deploymentParams) initChain() error {
	// import genesis configuration
	cc, err := chain.Import(p.genesisPath)
	if err != nil {
		return fmt.Errorf(
			"failed to load chain config from %s: %w",
			p.genesisPath,
			err,
		)
	}

	// set genesis configuration
	p.genesisConfig = cc

	return nil
}

func (p *deploymentParams) updateGenesisConfig() error {
	// Fetch contract deployment allowlist from genesis config
	deploymentAllowlist, err := config.GetDeploymentAllowlist(p.genesisConfig)
	if err != nil {
		return err
	}

	doesExist := map[types.Address]bool{}

	for _, a := range deploymentAllowlist {
		doesExist[a] = true
	}

	for _, a := range p.addAddresses {
		doesExist[a] = true
	}

	for _, a := range p.removeAddresses {
		doesExist[a] = false
	}

	newDeploymentAllowlist := make([]types.Address, 0)

	for addr, exists := range doesExist {
		if exists {
			newDeploymentAllowlist = append(newDeploymentAllowlist, addr)
		}
	}

	// Set allowlist in genesis configuration
	allowlistConfig := config.GetAllowlist(p.genesisConfig)

	if allowlistConfig == nil {
		allowlistConfig = &chain.Allowlists{}
	}

	allowlistConfig.Deployment = newDeploymentAllowlist
	p.genesisConfig.Params.Allowlists = allowlistConfig

	// Save allowlist for result
	p.allowlist = newDeploymentAllowlist

	return nil
}

func (p *deploymentParams) overrideGenesisConfig() error {
	// Remove the current genesis configuration from the disk
	if err := os.Remove(p.genesisPath); err != nil {
		return err
	}

	// Save the new genesis configuration
	if err := helper.WriteGenesisConfigToDisk(
		p.genesisConfig,
		p.genesisPath,
	); err != nil {
		return err
	}

	return nil
}

func (p *deploymentParams) getResult() command.CommandResult {
	result := &DeploymentResult{
		AddAddresses:    p.addAddresses,
		RemoveAddresses: p.removeAddresses,
		Allowlist:       p.allowlist,
	}

	return result
}

func unmarshallRawAddresses(addresses []string) []types.Address {
	marshalledAddresses := make([]types.Address, len(addresses))

	for indx, address := range addresses {
		marshalledAddresses[indx] = types.StringToAddress(address)
	}

	return marshalledAddresses
}
