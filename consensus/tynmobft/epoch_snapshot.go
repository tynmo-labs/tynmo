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

var activeEpochProposerSnapshot *EpochProposerSnapshot = nil

var (
	errHeightSyncIncomplete = errors.New("height syncing is not complete, tolerate and wait")
	errNoValidState         = errors.New("there is not valid local state")
	errUnexpectedMap        = errors.New("unexpected snapshot map: can not find height")
	errSyncSnapshotEmpty    = errors.New("snapshot sync result empty")
)

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

type SpsStatusType uint

const (
	SpsStatusInit     SpsStatusType = iota // init
	SpsStatusSynced                        // snapshot success synced from other peers
	SpsStatusLoaded                        // snapshot success load from db
	SpsStatusSyncNone                      // snapshot not in db and sync call success but result is nil, to calc
	SpsStatusCalced                        // snapshot calculated
)

type EpochProposerSnapshot struct {
	CurEpochHeightBase            uint64
	ProposerSnapshotMap           map[uint64]*ProposerSnapshot // height -> *ProposerSnapshot
	TotalVotingPower              *big.Int
	Logger                        hclog.Logger // Reference to the logging
	backendConsensus              *backendIBFT
	Status                        SpsStatusType
	PrioritizedValidatorAddresses []types.Address
}

func GetEpochProposerSnapshotResult(eps *EpochProposerSnapshot) *types.EpochProposerSnapshotResult {
	var epsr types.EpochProposerSnapshotResult
	epsr.CurEpochHeightBase = eps.CurEpochHeightBase
	for i := uint64(0); i < eps.backendConsensus.GetEpochSize(); i++ {
		curHeight := eps.CurEpochHeightBase + uint64(i)
		ps := eps.ProposerSnapshotMap[curHeight]
		epsr.PrioritizedValidatorAddresses = append(epsr.PrioritizedValidatorAddresses, ps.Proposer.Metadata.Address)
	}

	return &epsr
}

func (eps *EpochProposerSnapshot) TrimRoundProposerMap(height uint64) {
	curHeight := height
	for {
		nextHeight := curHeight + 1
		ps, ok := eps.ProposerSnapshotMap[nextHeight]
		if !ok {
			delete(eps.ProposerSnapshotMap, curHeight)
			break
		}
		eps.ProposerSnapshotMap[curHeight] = ps
		curHeight = nextHeight
	}
}

func (eps *EpochProposerSnapshot) cleanupHistorySnapshotMap(height uint64) {
	for i := uint64(0); i < eps.backendConsensus.GetEpochSize(); i++ {
		curHeight := height - 1 - i
		_, ok := eps.ProposerSnapshotMap[curHeight]
		if ok {
			delete(eps.ProposerSnapshotMap, curHeight)
		}
	}
}

func (eps *EpochProposerSnapshot) StoreEpochSnapshotResultToLocalState() error {
	return eps.backendConsensus.StoreEpochSnapshotResult(GetEpochProposerSnapshotResult(eps))
}

func (eps *EpochProposerSnapshot) LoadFromAddress() bool {
	return eps.Status == SpsStatusLoaded || eps.Status == SpsStatusSynced
}

func (eps *EpochProposerSnapshot) CanCalc() bool {
	return eps.Status != SpsStatusInit
}

func (eps *EpochProposerSnapshot) GetProposerAddress(height uint64, round uint64) (*types.Address, error) {
	isEpochStart := eps.backendConsensus.IsStartOfEpoch(height)
	var err error = nil
	addressIdx := (height + round) % eps.backendConsensus.GetEpochSize()

	// Early height syncing might not be complete yet, so check local state first and wait until syncing
	// is done, before that just returns an error.
	if eps.Status == SpsStatusInit {
		epochSnapshotResult, err := eps.backendConsensus.GetEpochSnapshotResult()
		if err != nil {
			return nil, err
		}
		if epochSnapshotResult != nil &&
			epochSnapshotResult.CurEpochHeightBase == eps.backendConsensus.GetEpochBaseHeight(height) &&
			len(epochSnapshotResult.PrioritizedValidatorAddresses) > 0 {
			eps.Status = SpsStatusLoaded
			eps.CurEpochHeightBase = epochSnapshotResult.CurEpochHeightBase
			eps.PrioritizedValidatorAddresses = epochSnapshotResult.PrioritizedValidatorAddresses
		} else {
			epochSnapshotResult, err = eps.backendConsensus.syncer.SyncEpochSnapshotOnce()
			if err != nil {
				return nil, err
			}
			if epochSnapshotResult == nil || len(epochSnapshotResult.PrioritizedValidatorAddresses) == 0 {
				eps.Status = SpsStatusSyncNone
			} else {
				eps.Status = SpsStatusSynced
				eps.CurEpochHeightBase = epochSnapshotResult.CurEpochHeightBase
				eps.PrioritizedValidatorAddresses = epochSnapshotResult.PrioritizedValidatorAddresses
			}
		}
	}

	if eps.LoadFromAddress() {
		if eps.CurEpochHeightBase >= eps.backendConsensus.GetEpochBaseHeight(height) {
			return &eps.PrioritizedValidatorAddresses[addressIdx], nil
		} else if !isEpochStart {
			eps.Status = SpsStatusInit
		}
	}

	// Calculate snapshot from stake contract
	if isEpochStart && eps.CanCalc() {
		if len(eps.ProposerSnapshotMap) == 0 || eps.CurEpochHeightBase != eps.backendConsensus.GetEpochBaseHeight(height) {
			eps.Logger.Info("SnapshotMap empty, CalculateAll", "height", height)
			err = eps.CalculateAll(height)
			if err != nil {
				eps.Logger.Error("CalculateAll error", "error", err)
				return nil, err
			}
			eps.Status = SpsStatusCalced
			// Save the snapshot to local state store
			eps.StoreEpochSnapshotResultToLocalState()
		}
	}

	snapshot, has := eps.ProposerSnapshotMap[addressIdx+eps.CurEpochHeightBase]
	if has {
		return &snapshot.Proposer.Metadata.Address, nil
	}
	return nil, errUnexpectedMap
}

func (eps *EpochProposerSnapshot) PreProposerSnapshot(height uint64) (*ProposerSnapshot, error) {
	if len(eps.ProposerSnapshotMap) == 0 {
		return nil, fmt.Errorf("snapshot map should has value")
	}
	heights := make([]uint64, 0)
	for height := range eps.ProposerSnapshotMap {
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
	return eps.ProposerSnapshotMap[pre], nil
}

func (eps *EpochProposerSnapshot) CalculateNextEpoch(height uint64) error {
	prePs, err := eps.PreProposerSnapshot(height)
	if err != nil {
		return err
	}

	epoch := eps.backendConsensus.GetEpochBaseHeight(height)

	eps.ProposerSnapshotMap = make(map[uint64]*ProposerSnapshot)
	eps.ProposerSnapshotMap[prePs.Height] = prePs

	psNew := prePs.Copy()
	for i := uint64(0); i < eps.backendConsensus.GetEpochSize(); i++ {
		nextHeight := epoch + i
		psNext, err := eps.NextHeight(psNew, nextHeight, eps.TotalVotingPower)
		if err != nil {
			return err
		}

		eps.CurEpochHeightBase = epoch
		eps.ProposerSnapshotMap[nextHeight] = psNext
		psNew = psNext
	}

	return nil
}

func (eps *EpochProposerSnapshot) CalculateSkip(height uint64, round uint64) error {
	epoch := eps.backendConsensus.GetEpochBaseHeight(height)
	preMap := eps.ProposerSnapshotMap
	eps.ProposerSnapshotMap = make(map[uint64]*ProposerSnapshot)

	for i := uint64(0); i < eps.backendConsensus.GetEpochSize(); i++ {
		nextHeight := epoch + i
		psNew := preMap[nextHeight]
		switch {
		case nextHeight < height:
			eps.ProposerSnapshotMap[nextHeight] = preMap[nextHeight]
		case nextHeight == height:
			for j := 0; j < int(round); j++ {
				psNext, err := eps.NextHeight(psNew.Copy(), nextHeight, eps.TotalVotingPower)
				if err != nil {
					return err
				}
				psNew = psNext
			}
			psNew.Round = round
			eps.ProposerSnapshotMap[nextHeight] = psNew
		case nextHeight > height:
			psNext, err := eps.NextHeight(psNew.Copy(), nextHeight, eps.TotalVotingPower)
			if err != nil {
				return err
			}
			psNext.Round = 0
			eps.ProposerSnapshotMap[nextHeight] = psNext
		}
	}
	return nil
}

func (eps *EpochProposerSnapshot) NextHeight(ps *ProposerSnapshot, height uint64, totalVotingPower *big.Int) (*ProposerSnapshot, error) {
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

func (eps *EpochProposerSnapshot) GenerateProposerSnapshot(height uint64) (*ProposerSnapshot, error) {
	accountSet, err := eps.backendConsensus.GetAccountSet(height)
	if err != nil {
		eps.Logger.Debug("EpochProposerSnapshot: Calculate: err", err)
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

// CalculateAll computing is based on the epoch base height
func (eps *EpochProposerSnapshot) CalculateAll(height uint64) error {
	eps.Logger.Debug("EpochProposerSnapshot: Calculate: starts")
	epoch := eps.backendConsensus.GetEpochBaseHeight(height)
	// Cleanup useless history snapshots
	eps.cleanupHistorySnapshotMap(epoch)

	// First round
	ps, err := eps.GenerateProposerSnapshot(height)
	if err != nil {
		return err
	}

	eps.TotalVotingPower = ps.GetTotalVotingPower()

	// ps.Proposer.ProposerPriority.Sub(ps.Proposer.ProposerPriority, totalVotingPower)
	eps.CurEpochHeightBase = epoch
	eps.ProposerSnapshotMap[epoch] = ps

	// Second round to EpochSize round
	psFormer := ps
	for i := uint64(1); i < eps.backendConsensus.GetEpochSize(); i++ {
		nextHeight := epoch + i
		psNew, err := eps.NextHeight(psFormer.Copy(), nextHeight, eps.TotalVotingPower)
		if err != nil {
			return err
		}
		eps.ProposerSnapshotMap[nextHeight] = psNew
		psFormer = psNew
	}
	eps.debugPrint()
	return nil
}

func (eps *EpochProposerSnapshot) debugPrint() {
	eps.Logger.Debug("EpochProposerSnapshot Print", "CurEpochHeightBase", eps.CurEpochHeightBase)
	heights := make([]uint64, 0)
	for height := range eps.ProposerSnapshotMap {
		heights = append(heights, height)
	}
	// ascending order
	sort.Slice(heights, func(i, j int) bool {
		return heights[i] < heights[j]
	})
	for _, height := range heights {
		ps := eps.ProposerSnapshotMap[height]
		eps.Logger.Debug("EpochProposerSnapshot Print", "Height", height, "Voted Proposer address", ps.Proposer.Metadata.Address.String())
		for _, validator := range ps.Validators {
			eps.Logger.Debug("EpochProposerSnapshot Print", "Validator: ProposerPriority", validator.ProposerPriority,
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
