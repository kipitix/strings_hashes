// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: hashcalc/api/hashcalc.proto

package grpchashcalc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HashCalcClient is the client API for HashCalc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HashCalcClient interface {
	Calc(ctx context.Context, opts ...grpc.CallOption) (HashCalc_CalcClient, error)
}

type hashCalcClient struct {
	cc grpc.ClientConnInterface
}

func NewHashCalcClient(cc grpc.ClientConnInterface) HashCalcClient {
	return &hashCalcClient{cc}
}

func (c *hashCalcClient) Calc(ctx context.Context, opts ...grpc.CallOption) (HashCalc_CalcClient, error) {
	stream, err := c.cc.NewStream(ctx, &HashCalc_ServiceDesc.Streams[0], "/grpchashcalc.HashCalc/Calc", opts...)
	if err != nil {
		return nil, err
	}
	x := &hashCalcCalcClient{stream}
	return x, nil
}

type HashCalc_CalcClient interface {
	Send(*InItem) error
	Recv() (*OutItem, error)
	grpc.ClientStream
}

type hashCalcCalcClient struct {
	grpc.ClientStream
}

func (x *hashCalcCalcClient) Send(m *InItem) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hashCalcCalcClient) Recv() (*OutItem, error) {
	m := new(OutItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HashCalcServer is the server API for HashCalc service.
// All implementations must embed UnimplementedHashCalcServer
// for forward compatibility
type HashCalcServer interface {
	Calc(HashCalc_CalcServer) error
	mustEmbedUnimplementedHashCalcServer()
}

// UnimplementedHashCalcServer must be embedded to have forward compatible implementations.
type UnimplementedHashCalcServer struct {
}

func (UnimplementedHashCalcServer) Calc(HashCalc_CalcServer) error {
	return status.Errorf(codes.Unimplemented, "method Calc not implemented")
}
func (UnimplementedHashCalcServer) mustEmbedUnimplementedHashCalcServer() {}

// UnsafeHashCalcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HashCalcServer will
// result in compilation errors.
type UnsafeHashCalcServer interface {
	mustEmbedUnimplementedHashCalcServer()
}

func RegisterHashCalcServer(s grpc.ServiceRegistrar, srv HashCalcServer) {
	s.RegisterService(&HashCalc_ServiceDesc, srv)
}

func _HashCalc_Calc_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HashCalcServer).Calc(&hashCalcCalcServer{stream})
}

type HashCalc_CalcServer interface {
	Send(*OutItem) error
	Recv() (*InItem, error)
	grpc.ServerStream
}

type hashCalcCalcServer struct {
	grpc.ServerStream
}

func (x *hashCalcCalcServer) Send(m *OutItem) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hashCalcCalcServer) Recv() (*InItem, error) {
	m := new(InItem)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HashCalc_ServiceDesc is the grpc.ServiceDesc for HashCalc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HashCalc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpchashcalc.HashCalc",
	HandlerType: (*HashCalcServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Calc",
			Handler:       _HashCalc_Calc_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hashcalc/api/hashcalc.proto",
}
