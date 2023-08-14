package syncer

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"tynmo/types"
)

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
