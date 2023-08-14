package syncer

import (
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"tynmo/types"
)

type SprintLocker struct {
	sprint uint64
	sync.RWMutex
}

// initNewPeerStatus fetches status of the peer and put to peer map
func (s *syncer) initNewPeerSnapshot(peerID peer.ID) {
	peer, err := s.syncPeerClient.GetSprintSnapshot(peerID, s.blockTimeout)
	if err != nil {
		s.logger.Warn("failed to get peer status, skip", "id", peerID, "err", err)
		return
	}

	s.peerSnapshotMap.PutSnapshots(peer)
	select {
	case s.newSnapshotCh <- struct{}{}:
	default:
	}
}

func (s *syncer) SyncSprintSnapshot(callback func(*types.SprintProposerSnapshotResult)) error {
	var wg sync.WaitGroup
	var wgCount = s.consensus.WaitPeerCount()
	var skipList = make(map[peer.ID]bool)

	wg.Add(wgCount)
	go func() {
		for i := 0; i < wgCount; i++ {
			<-s.newSnapshotCh
			localLatest := s.consensus.SprintHeightBase()
			bestPeer := s.peerSnapshotMap.BestSnapshotPeer(skipList)
			if bestPeer == nil {
				skipList = make(map[peer.ID]bool)
			} else {
				// if the bestPeer does not have a new block continue
				if bestPeer.Result.CurSprintHeightBase > localLatest {
					skipList[bestPeer.ID] = true
					callback(&bestPeer.Result)
				}
			}
			wg.Done()
		}
	}()

	wg.Wait()

	return nil
}

func (s *syncer) SyncSprintSnapshotOnce() (*types.SprintProposerSnapshotResult, error) {
	localLatest := s.blockchain.Header().Number
	bestPeer := s.peerMap.BestPeer(nil)

	if bestPeer == nil {
		return nil, fmt.Errorf("peer not ready")
	}

	if localLatest > bestPeer.Number {
		s.logger.Info("no need to sync snapshot", "peer.number", bestPeer.Number, "local", localLatest)
		return nil, nil
	}

	s.sprintLocker.Lock()
	defer s.sprintLocker.Unlock()

	sprint := s.consensus.GetSprint(bestPeer.Number)

	if sprint == s.sprintLocker.sprint {
		return nil, nil
	}

	if sprint != s.sprintLocker.sprint {
		s.sprintLocker.sprint = sprint
	}

	snapshot, err := s.syncPeerClient.GetSprintSnapshot(bestPeer.ID, s.blockTimeout)
	if err != nil {
		s.logger.Error("failed to get peer status, skip", "id", bestPeer.ID, "err", err)
		return nil, err
	}

	err = s.consensus.StoreSprintSnapshotResult(&snapshot.Result)
	if err != nil {
		return nil, err
	}
	return &snapshot.Result, nil
}
