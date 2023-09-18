package show

import (
	"crypto/ecdsa"
	"errors"
	"path/filepath"

	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/consensus/tynmobft"
	"tynmo/crypto"
	"tynmo/types"
)

const (
	AccountDirFlag       = "data-dir"
	PrivateKeyFlag       = "private-key"
	ValidatorAddressFlag = "address"
	OperateFlag          = "operate"

	AccountDirFlagDesc   = "the directory for the tynmo chain data if the local FS is used"
	ValidatorAddressDesc = "the account address to operate"
	OperateDesc          = "operate type: add/remove/get"
	PrivateKeyDesc       = "private key of the validator"
)

var (
	errNoNewValidatorsProvided          = errors.New("no new validators addresses provided")
	errPrivateKeyOrLocalDirNotSpecified = errors.New("only one of private-key and data-dir must be specified")
	errValidOperateType                 = errors.New("invalid operate type, only add/remove")
)

type whitelistParams struct {
	// private key related
	accountDir    string
	privateKeyStr string
	privateKey    *ecdsa.PrivateKey

	// json rpc url with http protocol by default
	jsonRPC string
}

func (p *whitelistParams) validateFlags() error {
	if (p.accountDir == "" && p.privateKeyStr == "") ||
		(p.accountDir != "" && p.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	_, err := helper.ParseJSONRPCAddress(p.jsonRPC)

	return err
}

func (p *whitelistParams) initRawParams() error {
	if err := p.initPrivateKey(); err != nil {
		return err
	}
	return nil
}

func (p *whitelistParams) initPrivateKey() error {
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
func (p *whitelistParams) initPrivateKeyFromLocalDataDir() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateOrReadPrivateKey(filepath.Join(p.accountDir, "consensus", tynmobft.IbftKeyName))
}

func (p *whitelistParams) getResult(list []string) command.CommandResult {
	return &showWhitelistResult{
		Validators: list,
	}
}
