// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.5
// source: consensus/proto/bft_operator.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BftOperator_GetSnapshot_FullMethodName = "/v1.BftOperator/GetSnapshot"
	BftOperator_Propose_FullMethodName     = "/v1.BftOperator/Propose"
	BftOperator_Candidates_FullMethodName  = "/v1.BftOperator/Candidates"
	BftOperator_Status_FullMethodName      = "/v1.BftOperator/Status"
)

// BftOperatorClient is the client API for BftOperator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BftOperatorClient interface {
	GetSnapshot(ctx context.Context, in *SnapshotReq, opts ...grpc.CallOption) (*Snapshot, error)
	Propose(ctx context.Context, in *Candidate, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Candidates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CandidatesResp, error)
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BftStatusResp, error)
}

type bftOperatorClient struct {
	cc grpc.ClientConnInterface
}

func NewBftOperatorClient(cc grpc.ClientConnInterface) BftOperatorClient {
	return &bftOperatorClient{cc}
}

func (c *bftOperatorClient) GetSnapshot(ctx context.Context, in *SnapshotReq, opts ...grpc.CallOption) (*Snapshot, error) {
	out := new(Snapshot)
	err := c.cc.Invoke(ctx, BftOperator_GetSnapshot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bftOperatorClient) Propose(ctx context.Context, in *Candidate, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BftOperator_Propose_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bftOperatorClient) Candidates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CandidatesResp, error) {
	out := new(CandidatesResp)
	err := c.cc.Invoke(ctx, BftOperator_Candidates_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bftOperatorClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BftStatusResp, error) {
	out := new(BftStatusResp)
	err := c.cc.Invoke(ctx, BftOperator_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BftOperatorServer is the server API for BftOperator service.
// All implementations must embed UnimplementedBftOperatorServer
// for forward compatibility
type BftOperatorServer interface {
	GetSnapshot(context.Context, *SnapshotReq) (*Snapshot, error)
	Propose(context.Context, *Candidate) (*emptypb.Empty, error)
	Candidates(context.Context, *emptypb.Empty) (*CandidatesResp, error)
	Status(context.Context, *emptypb.Empty) (*BftStatusResp, error)
	mustEmbedUnimplementedBftOperatorServer()
}

// UnimplementedBftOperatorServer must be embedded to have forward compatible implementations.
type UnimplementedBftOperatorServer struct {
}

func (UnimplementedBftOperatorServer) GetSnapshot(context.Context, *SnapshotReq) (*Snapshot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapshot not implemented")
}
func (UnimplementedBftOperatorServer) Propose(context.Context, *Candidate) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Propose not implemented")
}
func (UnimplementedBftOperatorServer) Candidates(context.Context, *emptypb.Empty) (*CandidatesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Candidates not implemented")
}
func (UnimplementedBftOperatorServer) Status(context.Context, *emptypb.Empty) (*BftStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedBftOperatorServer) mustEmbedUnimplementedBftOperatorServer() {}

// UnsafeBftOperatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BftOperatorServer will
// result in compilation errors.
type UnsafeBftOperatorServer interface {
	mustEmbedUnimplementedBftOperatorServer()
}

func RegisterBftOperatorServer(s grpc.ServiceRegistrar, srv BftOperatorServer) {
	s.RegisterService(&BftOperator_ServiceDesc, srv)
}

func _BftOperator_GetSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SnapshotReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BftOperatorServer).GetSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BftOperator_GetSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BftOperatorServer).GetSnapshot(ctx, req.(*SnapshotReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BftOperator_Propose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Candidate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BftOperatorServer).Propose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BftOperator_Propose_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BftOperatorServer).Propose(ctx, req.(*Candidate))
	}
	return interceptor(ctx, in, info, handler)
}

func _BftOperator_Candidates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BftOperatorServer).Candidates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BftOperator_Candidates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BftOperatorServer).Candidates(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BftOperator_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BftOperatorServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BftOperator_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BftOperatorServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BftOperator_ServiceDesc is the grpc.ServiceDesc for BftOperator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BftOperator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.BftOperator",
	HandlerType: (*BftOperatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSnapshot",
			Handler:    _BftOperator_GetSnapshot_Handler,
		},
		{
			MethodName: "Propose",
			Handler:    _BftOperator_Propose_Handler,
		},
		{
			MethodName: "Candidates",
			Handler:    _BftOperator_Candidates_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _BftOperator_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consensus/proto/bft_operator.proto",
}
