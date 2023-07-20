package add_validator

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"tynmo/command"
	"tynmo/command/helper"
	ibftOp "tynmo/consensus/proto"
	"tynmo/crypto"
	"tynmo/types"
)

const (
	voteFlag    = "vote"
	addressFlag = "addr"
	fromFlag    = "from"
	blsFlag     = "bls"
)

const (
	authVote = "auth"
	dropVote = "drop"
)

var (
	ErrFromPositive         = errors.New(`"from" must be positive number`)
	errInvalidVoteType      = errors.New("invalid vote type")
	errInvalidAddressFormat = errors.New("invalid address format")
)

var (
	params = &addValidatorParams{}
)

type addValidatorParams struct {
	addressRaw      string
	rawBLSPublicKey string
	fromRaw         string

	vote         string
	address      types.Address
	blsPublicKey []byte
	from         uint64
}

func (p *addValidatorParams) getRequiredFlags() []string {
	return []string{
		voteFlag,
		addressFlag,
		fromFlag,
	}
}

func (p *addValidatorParams) validateFlags() error {
	if !isValidVoteType(p.vote) {
		return errInvalidVoteType
	}

	return nil
}

func (p *addValidatorParams) initRawParams() error {
	if err := p.initAddress(); err != nil {
		return err
	}

	if err := p.initBLSPublicKey(); err != nil {
		return err
	}

	if err := p.initFrom(); err != nil {
		return err
	}

	return nil
}

func (p *addValidatorParams) initAddress() error {
	p.address = types.Address{}
	if err := p.address.UnmarshalText([]byte(p.addressRaw)); err != nil {
		return errInvalidAddressFormat
	}

	return nil
}

func (p *addValidatorParams) initFrom() error {
	from, err := types.ParseUint64orHex(&p.fromRaw)
	if err != nil {
		return fmt.Errorf("unable to parse from value, %w", err)
	}

	if from <= 0 {
		return ErrFromPositive
	}

	p.from = from

	return nil
}

func (p *addValidatorParams) initBLSPublicKey() error {
	if p.rawBLSPublicKey == "" {
		return nil
	}

	blsPubkeyBytes, err := hex.DecodeString(strings.TrimPrefix(p.rawBLSPublicKey, "0x"))
	if err != nil {
		return fmt.Errorf("failed to parse BLS Public Key: %w", err)
	}

	if _, err := crypto.UnmarshalBLSPublicKey(blsPubkeyBytes); err != nil {
		return err
	}

	p.blsPublicKey = blsPubkeyBytes

	return nil
}

func isValidVoteType(vote string) bool {
	return vote == authVote || vote == dropVote
}

func (p *addValidatorParams) AddValidatorCandidate(grpcAddress string) error {
	ibftClient, err := helper.GetIBFTOperatorClientConnection(grpcAddress)
	if err != nil {
		return err
	}

	if _, err := ibftClient.AddValidator(
		context.Background(),
		p.getCandidate(),
	); err != nil {
		return err
	}

	return nil
}

func (p *addValidatorParams) getCandidate() *ibftOp.Candidate {
	res := &ibftOp.Candidate{
		Address: p.address.String(),
		Auth:    p.vote == authVote,
		From:    p.from,
	}

	if p.blsPublicKey != nil {
		res.BlsPubkey = p.blsPublicKey
	}

	return res
}

func (p *addValidatorParams) getResult() command.CommandResult {
	return &IBFTAddValidatorResult{
		Address: p.address.String(),
		Vote:    p.vote,
	}
}
