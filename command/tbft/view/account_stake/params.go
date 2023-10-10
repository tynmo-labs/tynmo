package account_stake

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
	AccountFlag    = "account"

	AccountDirFlagDesc = "the directory for the tynmo chain data if the local FS is used"
	PrivateKeyDesc     = "private key of the validator"
	AccountDesc        = "account address"
)

var (
	errPrivateKeyOrLocalDirNotSpecified = errors.New("only one of private-key and data-dir must be specified")
)

type viewParams struct {
	// private key related
	accountDir    string
	privateKeyStr string
	privateKey    *ecdsa.PrivateKey

	account string

	// json rpc url with http protocol by default
	jsonRPC string
}

func (p *viewParams) validateFlags() error {
	if (p.accountDir == "" && p.privateKeyStr == "") ||
		(p.accountDir != "" && p.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	_, err := helper.ParseJSONRPCAddress(p.jsonRPC)

	return err
}

func (p *viewParams) initRawParams() error {
	if err := p.initPrivateKey(); err != nil {
		return err
	}
	return nil
}

func (p *viewParams) initPrivateKey() error {
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
func (p *viewParams) initPrivateKeyFromLocalDataDir() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateOrReadPrivateKey(filepath.Join(p.accountDir, "consensus", tynmobft.IbftKeyName))
}

func (p *viewParams) getResult(amount *big.Int) command.CommandResult {
	return &viewResult{
		Amount: amount,
	}
}
