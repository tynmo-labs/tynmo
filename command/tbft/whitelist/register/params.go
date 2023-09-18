package register

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"tynmo/command"
	"tynmo/command/helper"
	"tynmo/consensus/tynmobft"
	"tynmo/crypto"
	"tynmo/helper/hex"
	"tynmo/secrets"
	"tynmo/types"

	"github.com/coinbase/kryptology/pkg/signatures/bls/bls_sig"
)

const (
	AccountDirFlag       = "data-dir"
	PrivateKeyFlag       = "private-key"
	ValidatorAddressFlag = "address"
	AmountFlag           = "amount"

	AccountDirFlagDesc = "the directory for the tynmo chain data if the local FS is used"
	AmountToStakeDesc  = "amount to stake for a validator"
	PrivateKeyDesc     = "private key of the validator"
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
	blsPubKey     string
	privateKey    *ecdsa.PrivateKey

	// json rpc url with http protocol by default
	jsonRPC string

	// amount to stake
	amount      string
	amountValue *big.Int

	// returned transaction hash value
	hashRet string
}

func (p *whitelistParams) validateFlags() error {
	_, err := helper.ParseJSONRPCAddress(p.jsonRPC)
	if err != nil {
		return err
	}
	if (p.accountDir == "" && p.privateKeyStr == "") ||
		(p.accountDir != "" && p.privateKeyStr != "") {
		return errPrivateKeyOrLocalDirNotSpecified
	}

	if p.amountValue, err = helper.ParseAmount(p.amount); err != nil {
		return err
	}

	return nil
}

func (p *whitelistParams) initRawParams() error {
	if err := p.initPrivateKey(); err != nil {
		return err
	}
	if err := p.initBlsKey(); err != nil {
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

func (p *whitelistParams) initBlsKey() error {
	blsKey, err := p.initBlsKeyFromLocalDataDir()
	if err != nil {
		return err
	}

	blsBinary, err := crypto.BLSSecretKeyToPubkeyBytes(blsKey)
	if err != nil {
		return err
	}
	p.blsPubKey = hex.EncodeToHex(blsBinary)
	return nil
}

// PrivateKey returns a private key in data directory
func (p *whitelistParams) initPrivateKeyFromLocalDataDir() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateOrReadPrivateKey(filepath.Join(p.accountDir, "consensus", tynmobft.IbftKeyName))
}

func (p *whitelistParams) initBlsKeyFromLocalDataDir() (*bls_sig.SecretKey, error) {
	secret, err := os.ReadFile(filepath.Join(p.accountDir, "consensus", secrets.ValidatorBLSKeyLocal))
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read secret from disk %w",
			err,
		)
	}

	return crypto.BytesToBLSSecretKey(secret)
}

func (p *whitelistParams) getResult(hashRet string, keyHashRet string) command.CommandResult {
	addr := crypto.PubKeyToAddress(&p.privateKey.PublicKey)
	return &registerWhitelistResult{
		PublicAddress:     addr.String(),
		TxHashReturned:    hashRet,
		KeyTxHashReturned: keyHashRet,
	}
}
