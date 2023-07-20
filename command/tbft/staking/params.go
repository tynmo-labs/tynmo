package staking

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
	AmountFlag     = "amount"

	AccountDirFlagDesc = "the directory for the tynmo chain data if the local FS is used"
	AmountToStakeDesc  = "amount to stake for a validator"
	PrivateKeyDesc     = "private key of the validator"
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

	// amount to stake
	amount      string
	amountValue *big.Int

	// returned transaction hash value
	hashRet string
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
	if t.amountValue, err = helper.ParseAmount(t.amount); err != nil {
		return err
	}

	if (t.accountDir == "" && t.privateKeyStr == "") ||
		(t.accountDir != "" && t.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	return nil
}

func (t *stakeParams) getResult(hashRet string) command.CommandResult {
	addr := crypto.PubKeyToAddress(&t.privateKey.PublicKey)
	return &IBFTStakeResult{
		PublicAddress:  addr.String(),
		TxHashReturned: hashRet,
	}
}
