package tynmobft

import (
	"bytes"
	"fmt"
	"math/big"
)

// PrioritizedValidator holds ValidatorMetadata together with priority
type PrioritizedValidator struct {
	Metadata         *ValidatorMetadata
	ProposerPriority *big.Int
}

// ProposerSnapshot represents snapshot of one proposer calculation
type ProposerSnapshot struct {
	Height     uint64
	Proposer   *PrioritizedValidator
	Validators []*PrioritizedValidator
}

type SprintProposerSnapshot struct {
	CurSprintHeightBase uint64
	ProposerSnapshotMap map[uint64]*ProposerSnapshot // height -> *ProposerSnapshot
	backendConsensus    *backendIBFT
}

var activeSprintProposerSnapshot *SprintProposerSnapshot = nil

func (sps *SprintProposerSnapshot) GetProposerSnapshot(height uint64) *ProposerSnapshot {
	ps, ok := sps.ProposerSnapshotMap[height]
	if ok {
		return ps
	}

	err := sps.Calculate(height)
	if err != nil {
		return nil
	}

	ps, ok = sps.ProposerSnapshotMap[height]
	if ok {
		return ps
	}
	return nil
}

// Calculate computing is based on the sprint base height
func (sps *SprintProposerSnapshot) Calculate(height uint64) error {
	sps.backendConsensus.logger.Debug("SprintProposerSnapshot: Calculate: starts")
	accountSet, err := sps.backendConsensus.validatorsSnapshotCache.GetSnapshot(height)

	if err != nil {
		sps.backendConsensus.logger.Debug("SprintProposerSnapshot: Calculate: err", err)
		return err
	}

	sprint := GetSprint(height)
	ps := &ProposerSnapshot{
		Height:     sprint,
		Validators: make([]*PrioritizedValidator, 0, 0),
	}

	// First round
	for _, validator := range accountSet {
		pv := &PrioritizedValidator{
			Metadata: validator.Copy(),
			//ProposerPriority: big.NewInt(validator.VotingPower.Int64()),
			ProposerPriority: new(big.Int).Set(validator.VotingPower),
		}
		ps.Validators = append(ps.Validators, pv)
	}

	totalVotingPower := ps.GetTotalVotingPower()

	maxPrioritizedValidator, err := getValWithMostPriority(ps)
	if err != nil {
		return fmt.Errorf("cannot get validator with most priority for round 0: %w", err)
	}
	ps.Proposer = maxPrioritizedValidator

	ps.Proposer.ProposerPriority.Sub(ps.Proposer.ProposerPriority, totalVotingPower)
	sps.CurSprintHeightBase = sprint
	sps.ProposerSnapshotMap[sprint] = ps

	// Second round to SprintSize round
	psFormer := ps
	for i := 1; i < SprintSize; i++ {
		psNew := psFormer.Copy()
		for _, validator := range psNew.Validators {
			validator.ProposerPriority.Add(validator.ProposerPriority, validator.Metadata.VotingPower)
		}
		psNew.Height = sprint + uint64(i)
		newMaxPrioritizedValidator, err := getValWithMostPriority(psNew)
		if err != nil {
			return fmt.Errorf("cannot get validator with most priority for round %d: %w", i, err)
		}
		psNew.Proposer = newMaxPrioritizedValidator
		psNew.Proposer.ProposerPriority.Sub(psNew.Proposer.ProposerPriority, totalVotingPower)
		sps.CurSprintHeightBase = sprint
		sps.ProposerSnapshotMap[psNew.Height] = psNew

		psFormer = psNew
	}
	sps.debugPrint()
	return nil
}

func (sps *SprintProposerSnapshot) debugPrint() {
	sps.backendConsensus.logger.Debug("SprintProposerSnapshot Print", "CurSprintHeightBase", sps.CurSprintHeightBase)
	for height, ps := range sps.ProposerSnapshotMap {
		sps.backendConsensus.logger.Debug("SprintProposerSnapshot Print", "Height", height, "Voted Proposer address", ps.Proposer.Metadata.Address.String())
		for _, validator := range ps.Validators {
			sps.backendConsensus.logger.Debug("SprintProposerSnapshot Print", "Validator: ProposerPriority", validator.ProposerPriority,
				"validator addr", validator.Metadata.Address.String(),
				"validator vp", validator.Metadata.VotingPower,
			)
		}
	}
}

// GetTotalVotingPower returns total voting power from all the validators
func (pcs *ProposerSnapshot) GetTotalVotingPower() *big.Int {
	totalVotingPower := new(big.Int)
	for _, v := range pcs.Validators {
		totalVotingPower.Add(totalVotingPower, v.Metadata.VotingPower)
	}

	return totalVotingPower
}

// Copy Returns copy of current ProposerSnapshot object
func (pcs *ProposerSnapshot) Copy() *ProposerSnapshot {
	var proposer *PrioritizedValidator

	valCopy := make([]*PrioritizedValidator, len(pcs.Validators))

	for i, val := range pcs.Validators {
		valCopy[i] = &PrioritizedValidator{
			Metadata:         val.Metadata.Copy(),
			ProposerPriority: new(big.Int).Set(val.ProposerPriority)}

		if pcs.Proposer != nil && pcs.Proposer.Metadata.Address == val.Metadata.Address {
			proposer = valCopy[i]
		}
	}

	return &ProposerSnapshot{
		Validators: valCopy,
		Height:     pcs.Height,
		Proposer:   proposer,
	}
}

func getValWithMostPriority(snapshot *ProposerSnapshot) (result *PrioritizedValidator, err error) {
	if len(snapshot.Validators) == 0 {
		return nil, fmt.Errorf("validators cannot be nil or empty")
	}

	for _, curr := range snapshot.Validators {
		// pick curr as result if it has greater priority
		// or if it has same priority but "smaller" address
		if isBetterProposer(curr, result) {
			result = curr
		}
	}

	return result, nil
}

// isBetterProposer compares provided PrioritizedValidator instances
// and chooses either one with higher ProposerPriority or the one with the smaller address (compared lexicographically).
func isBetterProposer(a, b *PrioritizedValidator) bool {
	if b == nil || a.ProposerPriority.Cmp(b.ProposerPriority) > 0 {
		return true
	} else if a.ProposerPriority == b.ProposerPriority {
		return bytes.Compare(a.Metadata.Address.Bytes(), b.Metadata.Address.Bytes()) <= 0
	}

	return false
}
