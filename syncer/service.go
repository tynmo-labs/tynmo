package syncer

import (
	"context"
	"errors"

	"tynmo/network/grpc"
	"tynmo/syncer/proto"
	"tynmo/types"

	"github.com/golang/protobuf/ptypes/empty"
)

var (
	ErrBlockNotFound = errors.New("block not found")
)

type syncPeerService struct {
	proto.UnimplementedSyncPeerServer

	consensus  Consensus        // reference to the consensus module
	blockchain Blockchain       // reference to the blockchain module
	network    Network          // reference to the network module
	config     Config           // reference to the chain config
	stream     *grpc.GrpcStream // reference to the grpc stream
}

func NewSyncPeerService(
	network Network,
	blockchain Blockchain,
	consensus Consensus,
	config Config,
) SyncPeerService {
	return &syncPeerService{
		blockchain: blockchain,
		network:    network,
		consensus:  consensus,
		config:     config,
	}
}

// Start starts syncPeerService
func (s *syncPeerService) Start() {
	s.setupGRPCServer()
}

// Close closes syncPeerService
func (s *syncPeerService) Close() error {
	return s.stream.Close()
}

// setupGRPCServer setup GRPC server
func (s *syncPeerService) setupGRPCServer() {
	s.stream = grpc.NewGrpcStream()

	proto.RegisterSyncPeerServer(s.stream.GrpcServer(), s)
	s.stream.Serve()
	s.network.RegisterProtocol(syncerProto, s.stream)
}

// GetBlocks is a gRPC endpoint to return blocks from the specific height via stream
func (s *syncPeerService) GetBlocks(
	req *proto.GetBlocksRequest,
	stream proto.SyncPeer_GetBlocksServer,
) error {
	// from to latest
	for i := req.From; i <= s.blockchain.Header().Number; i++ {
		block, ok := s.blockchain.GetBlockByNumber(i, true)
		if !ok {
			return ErrBlockNotFound
		}

		resp := toProtoBlock(block)

		// if client closes stream, context.Canceled is given
		if err := stream.Send(resp); err != nil {
			break
		}
	}

	return nil
}

// GetStatus is a gRPC endpoint to return the latest block number as a node status
func (s *syncPeerService) GetStatus(
	ctx context.Context,
	req *empty.Empty,
) (*proto.SyncPeerStatus, error) {
	var number uint64

	if s.blockchain != nil {
		if header := s.blockchain.Header(); header != nil {
			number = header.Number
		}
	}

	return &proto.SyncPeerStatus{
		Number: number,
	}, nil
}

// toProtoBlock converts type.Block -> proto.Block
func toProtoBlock(block *types.Block) *proto.Block {
	return &proto.Block{
		Block: block.MarshalRLP(),
	}
}

// GetInitConfig is a gRPC endpoint to return the chain config
func (s *syncPeerService) GetInitConfig(
	ctx context.Context,
	req *proto.GetInitConfigRequest,
) (*proto.InitConfig, error) {
	content, err := s.config.Export()
	if err != nil {
		return nil, err
	}
	return &proto.InitConfig{
		Content: content,
	}, nil
}

// GetEpochSnapshot is a gRPC endpoint to return the latest epoch snapshot result
func (s *syncPeerService) GetEpochSnapshot(
	ctx context.Context,
	req *empty.Empty,
) (*proto.EpochSnapshot, error) {
	snapshotResult, err := s.consensus.GetEpochSnapshotResult()
	if err != nil {
		return nil, err
	}

	addresses := make([][]byte, 0)
	for _, snapshot := range snapshotResult.PrioritizedValidatorAddresses {
		addresses = append(addresses, snapshot.Bytes())
	}
	return &proto.EpochSnapshot{
		EpochHeight: snapshotResult.CurEpochHeightBase,
		Addresses:   addresses,
	}, nil
}
