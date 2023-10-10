package delegatee_epoch

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"path/filepath"

	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/consensus/tynmobft"
	"tynmo/crypto"
	"tynmo/types"
)

const (
	AccountDirFlag = "data-dir"
	PrivateKeyFlag = "private-key"
	EndEpochFlag   = "end-epoch"
	EpochIdFlag    = "epoch-id"

	AccountDirFlagDesc = "the directory for the tynmo chain data if the local FS is used"
	PrivateKeyDesc     = "private key of the validator"
	EndEpochDesc       = "end-epoch"
	EpochIdDesc        = "epoch-id"
)

var (
	errPrivateKeyOrLocalDirNotSpecified = errors.New("only one of private-key and data-dir must be specified")
)

type registerParams struct {
	// private key related
	accountDir    string
	privateKeyStr string
	privateKey    *ecdsa.PrivateKey

	endEpoch      string
	endEpochValue *big.Int

	epochId      string
	epochIdValue *big.Int

	// json rpc url with http protocol by default
	jsonRPC string
}

func (p *registerParams) validateFlags() error {
	var err error
	if (p.accountDir == "" && p.privateKeyStr == "") ||
		(p.accountDir != "" && p.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	if p.endEpochValue, err = helper.ParseAmount(p.endEpoch); err != nil {
		return err
	}

	if p.epochIdValue, err = helper.ParseAmount(p.epochId); err != nil {
		return err
	}

	_, err = helper.ParseJSONRPCAddress(p.jsonRPC)
	return err
}

func (p *registerParams) initRawParams() error {
	if err := p.initPrivateKey(); err != nil {
		return err
	}
	return nil
}

func (p *registerParams) initPrivateKey() error {
	var err error
	if p.privateKeyStr != "" {
		p.privateKey, err = crypto.ParseECDSAPrivateKey(types.StringToBytes(p.privateKeyStr))
		return err
	} else {
		p.privateKey, err = p.initPrivateKeyFromLocalDataDir()
		return err
	}
}

// PrivateKey returns a private key in data directory
func (p *registerParams) initPrivateKeyFromLocalDataDir() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateOrReadPrivateKey(filepath.Join(p.accountDir, "consensus", tynmobft.IbftKeyName))
}

func (p *registerParams) getResult(hashRet string) command.CommandResult {
	return &delegateeResult{
		EndEpoch: p.endEpochValue,
		// PublicAddress:  p.publicAddress,
		TxHashReturned: hashRet,
	}
}
