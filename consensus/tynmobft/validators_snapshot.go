package tynmobft

import (
	"fmt"
	"sync"

	"tynmo/types"

	"github.com/hashicorp/go-hclog"
)

const (
	// validatorSnapshotLimit defines a maximum number of validator snapshots
	// that can be stored in cache (both memory and db)
	validatorSnapshotLimit = 100
	// numberOfSnapshotsToLeaveInMemory defines a number of validator snapshots to leave in memory
	numberOfSnapshotsToLeaveInMemory = 12
)

type validatorSnapshot struct {
	Epoch    uint64           `json:"epoch"`
	Snapshot types.AccountSet `json:"snapshot"`
}

func (vs *validatorSnapshot) copy() *validatorSnapshot {
	copiedAccountSet := vs.Snapshot.Copy()

	return &validatorSnapshot{
		Epoch:    vs.Epoch,
		Snapshot: copiedAccountSet,
	}
}

type validatorsSnapshotCache struct {
	snapshots        map[uint64]*validatorSnapshot
	backendConsensus *backendIBFT
	lock             sync.Mutex
	logger           hclog.Logger
}

// newValidatorsSnapshotCache initializes a new instance of validatorsSnapshotCache
func newValidatorsSnapshotCache(logger hclog.Logger, backendConsensus *backendIBFT) *validatorsSnapshotCache {
	return &validatorsSnapshotCache{
		snapshots:        map[uint64]*validatorSnapshot{},
		logger:           logger.Named("validators_snapshot"),
		backendConsensus: backendConsensus,
	}
}

// GetSnapshot tries to retrieve the most recent cached snapshot (if any) and
// applies pending validator set deltas to it.
// Otherwise, it builds a snapshot from scratch and applies pending validator set deltas.
func (v *validatorsSnapshotCache) GetSnapshot(height uint64) (types.AccountSet, error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	epochToGetSnapshot := v.backendConsensus.GetEpochBaseHeight(height)
	v.logger.Trace("Retrieving snapshot started...", "Height:", height, "Epoch:", epochToGetSnapshot)

	latestValidatorSnapshot, err := v.getLastCachedSnapshot(epochToGetSnapshot)
	if err != nil {
		return nil, err
	}

	if latestValidatorSnapshot != nil {
		// we have snapshot for required block (epoch) in cache
		return latestValidatorSnapshot.Snapshot, nil
	}

	// latestValidatorSnapshot == nil
	// Haven't managed to retrieve snapshot for any epoch from the cache.
	// Build snapshot from the scratch, by applying delta from the genesis block.
	validatorSnapshot, err := v.computeSnapshot(height)
	if err != nil {
		return nil, fmt.Errorf("failed to compute snapshot for epoch 0: %w", err)
	}

	v.logger.Trace("TODO: Built validators snapshot for genesis block")

	if err := v.cleanup(); err != nil {
		// error on cleanup should not block or fail any action
		v.logger.Error("could not clean validator snapshots from cache and db", "err", err)
	}

	return validatorSnapshot.Snapshot, nil
}

// computeSnapshot gets desired block header by block number, extracts its extra and applies given delta to the snapshot
func (v *validatorsSnapshotCache) computeSnapshot(height uint64) (*validatorSnapshot, error) {
	validators, err := v.backendConsensus.forkManager.GetValidators(height)
	if err != nil {
		return nil, err
	}

	v.logger.Trace("Compute snapshot started...", "Height", height)

	vs := validatorSnapshot{
		Epoch:    v.backendConsensus.GetEpochBaseHeight(height),
		Snapshot: types.NewAccountSet(),
	}

	for idx := 0; idx < validators.Len(); idx++ {
		validator := validators.At(uint64(idx))

		v.logger.Debug("validatorsSnapshotCache: computeSnapshot:", "Validator: Type",
			validator.Type(), "Address", validator.Addr().String(),
			"Stake", validator.GetStake())
		types.AppendAccountSet(&vs.Snapshot, types.NewValidatorMetadata(
			validator.Addr(), validator.GetStake(), true, true))
	}

	// Insert validatorSnapshot per epoch only
	v.snapshots[v.backendConsensus.GetEpochBaseHeight(height)] = &vs
	return &vs, nil
}

// Cleanup cleans the validators cache in memory and db
func (v *validatorsSnapshotCache) cleanup() error {
	if len(v.snapshots) >= validatorSnapshotLimit {
		latestEpoch := uint64(0)

		for e := range v.snapshots {
			if e > latestEpoch {
				latestEpoch = e
			}
		}

		startEpoch := latestEpoch
		cache := make(map[uint64]*validatorSnapshot, numberOfSnapshotsToLeaveInMemory)

		for i := 0; i < numberOfSnapshotsToLeaveInMemory; i++ {
			if snapshot, exists := v.snapshots[startEpoch]; exists {
				cache[startEpoch] = snapshot
			}

			startEpoch -= v.backendConsensus.GetEpochSize()
		}

		v.snapshots = cache

		return nil
	}

	return nil
}

// getLastCachedSnapshot gets the latest snapshot cached
func (v *validatorsSnapshotCache) getLastCachedSnapshot(currentEpoch uint64) (*validatorSnapshot, error) {
	cachedSnapshot := v.snapshots[currentEpoch]
	return cachedSnapshot, nil
}
