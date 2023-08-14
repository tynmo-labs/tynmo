package tynmobft

import (
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"tynmo/types"
)

/*
Bolt DB schema:

proposer snapshot/
|--> proposerSnapshotKey - only current one snapshot is preserved -> *ProposerSnapshot (json marshalled)
*/
var (
	// bucket to store proposer calculator snapshot
	proposerSnapshotBucket = []byte("proposerSnapshot")
	// proposerSnapshotKey is a static key which is used to save latest proposer snapshot.
	// (there will always be one object in bucket)
	proposerSnapshotKey       = []byte("proposerSnapshotKey")
	errReadLocalSnapshotState = errors.New("failed to read local snapshot state")
	// bucket to store validator snapshots
	validatorSnapshotsBucket = []byte("validatorSnapshots")
)

type ProposerSnapshotStore struct {
	db *bolt.DB
}

// initialize creates necessary buckets in DB if they don't already exist
func (s *ProposerSnapshotStore) initialize(tx *bolt.Tx) error {
	if _, err := tx.CreateBucketIfNotExists(proposerSnapshotBucket); err != nil {
		return fmt.Errorf("failed to create bucket=%s: %w", string(validatorSnapshotsBucket), err)
	}

	return nil
}

// GetProposerSnapshotResult gets the latest sprint proposer snapshot
func (s *ProposerSnapshotStore) GetProposerSnapshotResult() (*types.SprintProposerSnapshotResult, error) {
	var snapshot *types.SprintProposerSnapshotResult

	err := s.db.View(func(tx *bolt.Tx) error {
		value := tx.Bucket(proposerSnapshotBucket).Get(proposerSnapshotKey)
		if value == nil {
			return nil
		}

		return json.Unmarshal(value, &snapshot)
	})

	if snapshot == nil {
		/*
			1 TODO should return error ?
			2 return empty obj to deal with initial chain
		*/
		snapshot = &types.SprintProposerSnapshotResult{}
	}
	return snapshot, err
}

// WriteProposerSnapshotResult writes the latest sprint proposer snapshot
func (s *ProposerSnapshotStore) WriteProposerSnapshotResult(snapshot *types.SprintProposerSnapshotResult) error {
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		// Cleanup the key / value first
		tx.Bucket(proposerSnapshotBucket).Delete(proposerSnapshotKey)
		return tx.Bucket(proposerSnapshotBucket).Put(proposerSnapshotKey, raw)
	})
}
