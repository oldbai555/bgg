// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: lb.proto

package lb

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LbClient is the client API for Lb service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LbClient interface {
}

type lbClient struct {
	cc grpc.ClientConnInterface
}

func NewLbClient(cc grpc.ClientConnInterface) LbClient {
	return &lbClient{cc}
}

// LbServer is the server API for Lb service.
// All implementations must embed UnimplementedLbServer
// for forward compatibility
type LbServer interface {
	mustEmbedUnimplementedLbServer()
}

// UnimplementedLbServer must be embedded to have forward compatible implementations.
type UnimplementedLbServer struct {
}

func (UnimplementedLbServer) mustEmbedUnimplementedLbServer() {}

// UnsafeLbServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LbServer will
// result in compilation errors.
type UnsafeLbServer interface {
	mustEmbedUnimplementedLbServer()
}

func RegisterLbServer(s grpc.ServiceRegistrar, srv LbServer) {
	s.RegisterService(&Lb_ServiceDesc, srv)
}

// Lb_ServiceDesc is the grpc.ServiceDesc for Lb service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lb_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lb.lb",
	HandlerType: (*LbServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "lb.proto",
}
