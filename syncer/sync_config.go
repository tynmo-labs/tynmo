package syncer

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/peer"
	"tynmo/syncer/proto"
)

func NewInitConfigClient(logger hclog.Logger, network Network, config Config) SyncPeerClient {
	return &syncPeerClient{
		logger:  logger.Named(SyncConfigClientLoggerName),
		network: network,
		config:  config,
		id:      network.AddrInfo().ID.String(),
	}
}

func (m *syncPeerClient) GetInitConfig(peerID peer.ID) ([]byte, error) {
	service := &syncPeerService{network: m.network}
	service.Start()

	conn, err := m.network.NewProtoConnection(syncerProto, peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to open a stream, err %w", err)
	}

	clt := proto.NewSyncPeerClient(conn)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), defaultTimeoutForStatus)
	defer cancel()

	resp, err := clt.GetInitConfig(timeoutCtx, &proto.GetInitConfigRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Content, nil
}
