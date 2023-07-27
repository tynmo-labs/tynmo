package server

import (
	"tynmo/consensus"
	consensusDev "tynmo/consensus/dev"
	consensusDummy "tynmo/consensus/dummy"
	consensusIBFT "tynmo/consensus/ibft"
	consensusTynmoBFT "tynmo/consensus/tynmobft"
	"tynmo/secrets"
	"tynmo/secrets/awsssm"
	"tynmo/secrets/gcpssm"
	"tynmo/secrets/hashicorpvault"
	"tynmo/secrets/local"
)

type ConsensusType string

const (
	DevConsensus      ConsensusType = "dev"
	IBFTConsensus     ConsensusType = "ibft"
	TynmoBFTConsensus ConsensusType = "tynmobft"
	DummyConsensus    ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:      consensusDev.Factory,
	IBFTConsensus:     consensusIBFT.Factory,
	TynmoBFTConsensus: consensusTynmoBFT.Factory,
	DummyConsensus:    consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
