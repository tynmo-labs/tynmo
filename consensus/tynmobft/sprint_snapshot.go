package tynmobft

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/hashicorp/go-hclog"
	"tynmo/types"
)

const SprintSize = 10

var activeSprintProposerSnapshot *SprintProposerSnapshot = nil

var (
	errHeightSyncIncomplete = errors.New("height syncing is not complete, tolerate and wait")
	errNoValidState         = errors.New("there is not valid local state")
)

func GetSprint(height uint64) uint64 {
	return height - height%SprintSize
}

func GetSprintRound(height uint64) uint64 {
	return height % SprintSize
}

func IsSprintStart(height uint64) bool {
	return height == 1 || GetSprintRound(height) == 0
}

// PrioritizedValidator holds ValidatorMetadata together with priority
type PrioritizedValidator struct {
	Metadata         *types.ValidatorMetadata
	ProposerPriority *big.Int
}

// ProposerSnapshot represents snapshot of one proposer calculation
type ProposerSnapshot struct {
	Round      uint64
	Height     uint64
	Proposer   *PrioritizedValidator
	Validators []*PrioritizedValidator
}

type SprintProposerSnapshot struct {
	CurSprintHeightBase uint64
	ProposerSnapshotMap map[uint64]*ProposerSnapshot // height -> *ProposerSnapshot
	TotalVotingPower    *big.Int
	Logger              hclog.Logger // Reference to the logging
	backendConsensus    *backendIBFT
}

func GetSprintProposerSnapshotResult(sps *SprintProposerSnapshot) *types.SprintProposerSnapshotResult {
	var spsr types.SprintProposerSnapshotResult
	spsr.CurSprintHeightBase = sps.CurSprintHeightBase
	for i := 0; i < SprintSize; i++ {
		curHeight := sps.CurSprintHeightBase + uint64(i)
		ps := sps.ProposerSnapshotMap[curHeight]
		spsr.PrioritizedValidatorAddresses = append(spsr.PrioritizedValidatorAddresses, ps.Proposer.Metadata.Address)
	}

	return &spsr
}

func (sps *SprintProposerSnapshot) TrimRoundProposerMap(height uint64) {
	curHeight := height
	for {
		nextHeight := curHeight + 1
		ps, ok := sps.ProposerSnapshotMap[nextHeight]
		if !ok {
			delete(sps.ProposerSnapshotMap, curHeight)
			break
		}
		sps.ProposerSnapshotMap[curHeight] = ps
		curHeight = nextHeight
	}
}

func (sps *SprintProposerSnapshot) cleanupHistorySnapshotMap(height uint64) {
	for i := uint64(0); i < SprintSize; i++ {
		curHeight := height - 1 - i
		_, ok := sps.ProposerSnapshotMap[curHeight]
		if ok {
			delete(sps.ProposerSnapshotMap, curHeight)
		}
	}
}

func (sps *SprintProposerSnapshot) StoreSprintSnapshotResultToLocalState() error {
	return sps.backendConsensus.StoreSprintSnapshotResult(GetSprintProposerSnapshotResult(sps))
}

func (sps *SprintProposerSnapshot) GetProposerAddress(height uint64, round uint64) (*types.Address, error) {
	isSprintStart := IsSprintStart(height)
	var err error = nil
	addressIdx := (height + round) % SprintSize

	// Early height syncing might not be complete yet, so check local state first and wait until syncing
	// is done, before that just returns an error.
	sprintSnapshotResult, err := sps.backendConsensus.GetSprintSnapshotResult()
	if err == nil && sprintSnapshotResult != nil && sprintSnapshotResult.CurSprintHeightBase > height {
		return nil, errHeightSyncIncomplete
	}

	// Calculate snapshot from stake contract
	if isSprintStart {
		if len(sps.ProposerSnapshotMap) == 0 || sps.CurSprintHeightBase != GetSprint(height) {
			sps.Logger.Debug("SnapshotMap empty, CalculateAll", "height", height)
			err = sps.CalculateAll(height)
			if err != nil {
				sps.Logger.Error("CalculateAll error", "error", err)
				return nil, err
			}
			// Save the snapshot to local state store
			sps.StoreSprintSnapshotResultToLocalState()
		}
		return &sps.ProposerSnapshotMap[addressIdx+sps.CurSprintHeightBase].Proposer.Metadata.Address, nil
	}

	// Read sprint snapshot from local state storage
	if err == nil && sprintSnapshotResult != nil && sprintSnapshotResult.CurSprintHeightBase == GetSprint(height) {
		sps.Logger.Debug("Succeeded to read validator addresses from local state")
		return &sprintSnapshotResult.PrioritizedValidatorAddresses[addressIdx], nil
	} else {
		// Directly return the error and wait for the state syncing from peers
		// If syncing succeeds, sprintSnapshotResult will contain for next few rounds
		sps.Logger.Debug("Failed to read validator addresses from local state")
		return nil, errNoValidState
	}
}

func (sps *SprintProposerSnapshot) PreProposerSnapshot(height uint64) (*ProposerSnapshot, error) {
	if len(sps.ProposerSnapshotMap) == 0 {
		return nil, fmt.Errorf("snapshot map should has value")
	}
	heights := make([]uint64, 0)
	for height := range sps.ProposerSnapshotMap {
		heights = append(heights, height)
	}
	//desc
	sort.Slice(heights, func(i, j int) bool {
		return heights[i] > heights[j]
	})
	var pre uint64
	for i := 0; i < len(heights); i++ {
		if heights[i] <= height {
			pre = heights[i]
			break
		}
	}
	return sps.ProposerSnapshotMap[pre], nil
}

func (sps *SprintProposerSnapshot) CalculateNextSprint(height uint64) error {
	prePs, err := sps.PreProposerSnapshot(height)
	if err != nil {
		return err
	}

	sprint := GetSprint(height)

	sps.ProposerSnapshotMap = make(map[uint64]*ProposerSnapshot)
	sps.ProposerSnapshotMap[prePs.Height] = prePs

	psNew := prePs.Copy()
	for i := 0; i < SprintSize; i++ {
		nextHeight := sprint + uint64(i)
		psNext, err := sps.NextHeight(psNew, nextHeight, sps.TotalVotingPower)
		if err != nil {
			return err
		}

		sps.CurSprintHeightBase = sprint
		sps.ProposerSnapshotMap[nextHeight] = psNext
		psNew = psNext
	}

	return nil
}

func (sps *SprintProposerSnapshot) CalculateSkip(height uint64, round uint64) error {
	sprint := GetSprint(height)
	preMap := sps.ProposerSnapshotMap
	sps.ProposerSnapshotMap = make(map[uint64]*ProposerSnapshot)

	for i := 0; i < SprintSize; i++ {
		nextHeight := sprint + uint64(i)
		psNew := preMap[nextHeight]
		switch {
		case nextHeight < height:
			sps.ProposerSnapshotMap[nextHeight] = preMap[nextHeight]
		case nextHeight == height:
			for j := 0; j < int(round); j++ {
				psNext, err := sps.NextHeight(psNew.Copy(), nextHeight, sps.TotalVotingPower)
				if err != nil {
					return err
				}
				psNew = psNext
			}
			psNew.Round = round
			sps.ProposerSnapshotMap[nextHeight] = psNew
		case nextHeight > height:
			psNext, err := sps.NextHeight(psNew.Copy(), nextHeight, sps.TotalVotingPower)
			if err != nil {
				return err
			}
			psNext.Round = 0
			sps.ProposerSnapshotMap[nextHeight] = psNext
		}
	}
	return nil
}

func (sps *SprintProposerSnapshot) NextHeight(ps *ProposerSnapshot, height uint64, totalVotingPower *big.Int) (*ProposerSnapshot, error) {
	psNew := ps.Copy()
	for _, validator := range psNew.Validators {
		validator.ProposerPriority.Add(validator.ProposerPriority, validator.Metadata.VotingPower)
	}
	psNew.Height = height
	newMaxPrioritizedValidator, err := getValWithMostPriority(psNew)
	if err != nil {
		return nil, fmt.Errorf("cannot get validator with most priority for : %v", err)
	}
	psNew.Proposer = newMaxPrioritizedValidator
	psNew.Proposer.ProposerPriority.Sub(psNew.Proposer.ProposerPriority, totalVotingPower)
	psNew.Round = 0

	// if err := updateWithChangeSet(psNew, totalVotingPower); err != nil {
	// 	return nil, err
	// }
	return psNew, nil
}

func (sps *SprintProposerSnapshot) GenerateProposerSnapshot(height uint64) (*ProposerSnapshot, error) {
	accountSet, err := sps.backendConsensus.GetAccountSet(height)
	if err != nil {
		sps.Logger.Debug("SprintProposerSnapshot: Calculate: err", err)
		return nil, err
	}

	ps := &ProposerSnapshot{
		Height:     height,
		Validators: make([]*PrioritizedValidator, 0, 0),
	}

	for _, validator := range accountSet {
		pv := &PrioritizedValidator{
			Metadata:         validator.Copy(),
			ProposerPriority: new(big.Int).Set(validator.VotingPower),
		}
		ps.Validators = append(ps.Validators, pv)
	}

	maxPrioritizedValidator, err := getValWithMostPriority(ps)
	if err != nil {
		return nil, fmt.Errorf("cannot get validator with most priority for round 0: %w", err)
	}
	ps.Proposer = maxPrioritizedValidator

	return ps, nil
}

// CalculateAll computing is based on the sprint base height
func (sps *SprintProposerSnapshot) CalculateAll(height uint64) error {
	sps.Logger.Debug("SprintProposerSnapshot: Calculate: starts")
	sprint := GetSprint(height)
	// Cleanup useless history snapshots
	sps.cleanupHistorySnapshotMap(sprint)

	// First round
	ps, err := sps.GenerateProposerSnapshot(height)
	if err != nil {
		return err
	}

	sps.TotalVotingPower = ps.GetTotalVotingPower()

	// ps.Proposer.ProposerPriority.Sub(ps.Proposer.ProposerPriority, totalVotingPower)
	sps.CurSprintHeightBase = sprint
	sps.ProposerSnapshotMap[sprint] = ps

	// Second round to SprintSize round
	psFormer := ps
	for i := 1; i < SprintSize; i++ {
		nextHeight := sprint + uint64(i)
		psNew, err := sps.NextHeight(psFormer.Copy(), nextHeight, sps.TotalVotingPower)
		if err != nil {
			return err
		}
		sps.ProposerSnapshotMap[nextHeight] = psNew
		psFormer = psNew
	}
	sps.debugPrint()
	return nil
}

func (sps *SprintProposerSnapshot) debugPrint() {
	sps.Logger.Debug("SprintProposerSnapshot Print", "CurSprintHeightBase", sps.CurSprintHeightBase)
	heights := make([]uint64, 0)
	for height := range sps.ProposerSnapshotMap {
		heights = append(heights, height)
	}
	// ascending order
	sort.Slice(heights, func(i, j int) bool {
		return heights[i] < heights[j]
	})
	for _, height := range heights {
		ps := sps.ProposerSnapshotMap[height]
		sps.Logger.Debug("SprintProposerSnapshot Print", "Height", height, "Voted Proposer address", ps.Proposer.Metadata.Address.String())
		for _, validator := range ps.Validators {
			sps.Logger.Debug("SprintProposerSnapshot Print", "Validator: ProposerPriority", validator.ProposerPriority,
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
