// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.5
// source: consensus/ibft/proto/ibft_operator.proto

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
	V0IbftOperator_GetSnapshot_FullMethodName = "/v1.V0IbftOperator/GetSnapshot"
	V0IbftOperator_Propose_FullMethodName     = "/v1.V0IbftOperator/Propose"
	V0IbftOperator_Candidates_FullMethodName  = "/v1.V0IbftOperator/Candidates"
	V0IbftOperator_Status_FullMethodName      = "/v1.V0IbftOperator/Status"
)

// V0IbftOperatorClient is the client API for V0IbftOperator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type V0IbftOperatorClient interface {
	GetSnapshot(ctx context.Context, in *V0SnapshotReq, opts ...grpc.CallOption) (*V0Snapshot, error)
	Propose(ctx context.Context, in *V0Candidate, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Candidates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*V0CandidatesResp, error)
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*V0IbftStatusResp, error)
}

type v0IbftOperatorClient struct {
	cc grpc.ClientConnInterface
}

func NewV0IbftOperatorClient(cc grpc.ClientConnInterface) V0IbftOperatorClient {
	return &v0IbftOperatorClient{cc}
}

func (c *v0IbftOperatorClient) GetSnapshot(ctx context.Context, in *V0SnapshotReq, opts ...grpc.CallOption) (*V0Snapshot, error) {
	out := new(V0Snapshot)
	err := c.cc.Invoke(ctx, V0IbftOperator_GetSnapshot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v0IbftOperatorClient) Propose(ctx context.Context, in *V0Candidate, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, V0IbftOperator_Propose_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v0IbftOperatorClient) Candidates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*V0CandidatesResp, error) {
	out := new(V0CandidatesResp)
	err := c.cc.Invoke(ctx, V0IbftOperator_Candidates_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v0IbftOperatorClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*V0IbftStatusResp, error) {
	out := new(V0IbftStatusResp)
	err := c.cc.Invoke(ctx, V0IbftOperator_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// V0IbftOperatorServer is the server API for V0IbftOperator service.
// All implementations must embed UnimplementedV0IbftOperatorServer
// for forward compatibility
type V0IbftOperatorServer interface {
	GetSnapshot(context.Context, *V0SnapshotReq) (*V0Snapshot, error)
	Propose(context.Context, *V0Candidate) (*emptypb.Empty, error)
	Candidates(context.Context, *emptypb.Empty) (*V0CandidatesResp, error)
	Status(context.Context, *emptypb.Empty) (*V0IbftStatusResp, error)
	mustEmbedUnimplementedV0IbftOperatorServer()
}

// UnimplementedV0IbftOperatorServer must be embedded to have forward compatible implementations.
type UnimplementedV0IbftOperatorServer struct {
}

func (UnimplementedV0IbftOperatorServer) GetSnapshot(context.Context, *V0SnapshotReq) (*V0Snapshot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapshot not implemented")
}
func (UnimplementedV0IbftOperatorServer) Propose(context.Context, *V0Candidate) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Propose not implemented")
}
func (UnimplementedV0IbftOperatorServer) Candidates(context.Context, *emptypb.Empty) (*V0CandidatesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Candidates not implemented")
}
func (UnimplementedV0IbftOperatorServer) Status(context.Context, *emptypb.Empty) (*V0IbftStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedV0IbftOperatorServer) mustEmbedUnimplementedV0IbftOperatorServer() {}

// UnsafeV0IbftOperatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to V0IbftOperatorServer will
// result in compilation errors.
type UnsafeV0IbftOperatorServer interface {
	mustEmbedUnimplementedV0IbftOperatorServer()
}

func RegisterV0IbftOperatorServer(s grpc.ServiceRegistrar, srv V0IbftOperatorServer) {
	s.RegisterService(&V0IbftOperator_ServiceDesc, srv)
}

func _V0IbftOperator_GetSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(V0SnapshotReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V0IbftOperatorServer).GetSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V0IbftOperator_GetSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V0IbftOperatorServer).GetSnapshot(ctx, req.(*V0SnapshotReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _V0IbftOperator_Propose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(V0Candidate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V0IbftOperatorServer).Propose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V0IbftOperator_Propose_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V0IbftOperatorServer).Propose(ctx, req.(*V0Candidate))
	}
	return interceptor(ctx, in, info, handler)
}

func _V0IbftOperator_Candidates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V0IbftOperatorServer).Candidates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V0IbftOperator_Candidates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V0IbftOperatorServer).Candidates(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _V0IbftOperator_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V0IbftOperatorServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V0IbftOperator_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V0IbftOperatorServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// V0IbftOperator_ServiceDesc is the grpc.ServiceDesc for V0IbftOperator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var V0IbftOperator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.V0IbftOperator",
	HandlerType: (*V0IbftOperatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSnapshot",
			Handler:    _V0IbftOperator_GetSnapshot_Handler,
		},
		{
			MethodName: "Propose",
			Handler:    _V0IbftOperator_Propose_Handler,
		},
		{
			MethodName: "Candidates",
			Handler:    _V0IbftOperator_Candidates_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _V0IbftOperator_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consensus/ibft/proto/ibft_operator.proto",
}
