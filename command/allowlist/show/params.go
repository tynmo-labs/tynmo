package show

import (
	"fmt"

	"tynmo/chain"
	"tynmo/command"
	"tynmo/helper/config"
	"tynmo/types"
)

const (
	chainFlag = "chain"
)

var (
	params = &showParams{}
)

type showParams struct {
	// genesis file path
	genesisPath string

	// deployment allowlist
	allowlists Allowlists
}

type Allowlists struct {
	deployment []types.Address
}

func (p *showParams) initRawParams() error {
	// init genesis configuration
	if err := p.initAllowlists(); err != nil {
		return err
	}

	return nil
}

func (p *showParams) initAllowlists() error {
	// import genesis configuration
	genesisConfig, err := chain.Import(p.genesisPath)
	if err != nil {
		return fmt.Errorf(
			"failed to load chain config from %s: %w",
			p.genesisPath,
			err,
		)
	}

	// fetch allowlists
	deploymentAllowlist, err := config.GetDeploymentAllowlist(genesisConfig)
	if err != nil {
		return err
	}

	// set allowlists
	p.allowlists = Allowlists{
		deployment: deploymentAllowlist,
	}

	return nil
}

func (p *showParams) getResult() command.CommandResult {
	result := &ShowResult{
		Allowlists: p.allowlists,
	}

	return result
}
