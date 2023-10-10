package unstake_delegator

import (
	"crypto/ecdsa"
	"errors"
	"path/filepath"

	"tynmo/command"
	"tynmo/consensus/tynmobft"
	"tynmo/crypto"
	"tynmo/types"
)

const (
	AccountDirFlag = "data-dir"
	PrivateKeyFlag = "private-key"
	AmountFlag     = "amount"
	DelegateeFlag  = "delegatee"

	AccountDirFlagDesc = "the directory for the tynmo chain data if the local FS is used"
	AmountToStakeDesc  = "amount to stake for a validator"
	PrivateKeyDesc     = "private key of the validator"
	DelegateeDesc      = "delegatee address"
)

var (
	errPrivateKeyOrLocalDirNotSpecified = errors.New("only one of private-key and data-dir must be specified")
)

type stakeParams struct {
	// private key related
	accountDir    string
	privateKeyStr string
	privateKey    *ecdsa.PrivateKey

	// json rpc url with http protocol by default
	jsonRPC string

	delegatee string
}

func (t *stakeParams) initRawParams() error {
	if err := t.initPrivateKey(); err != nil {
		return err
	}
	return nil
}

func (t *stakeParams) initPrivateKey() error {
	var err error
	if t.privateKeyStr != "" {
		t.privateKey, err = crypto.ParseECDSAPrivateKey(types.StringToBytes(t.privateKeyStr))
		return err
	} else {
		t.privateKey, err = t.initPrivateKeyFromLocalDataDir()
		return err
	}
}

// PrivateKey returns a private key in data directory
func (t *stakeParams) initPrivateKeyFromLocalDataDir() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateOrReadPrivateKey(filepath.Join(t.accountDir, "consensus", tynmobft.IbftKeyName))
}

func (t *stakeParams) validateFlags() (err error) {
	if (t.accountDir == "" && t.privateKeyStr == "") ||
		(t.accountDir != "" && t.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	return nil
}

func (t *stakeParams) getResult(hashRet string) command.CommandResult {
	addr := crypto.PubKeyToAddress(&t.privateKey.PublicKey)
	return &StakeResult{
		DelegateeAddress: t.delegatee,
		PublicAddress:    addr.String(),
		TxHashReturned:   hashRet,
	}
}
