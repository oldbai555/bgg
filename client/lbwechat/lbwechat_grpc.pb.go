// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: lbwechat.proto

package lbwechat

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

// LbwechatClient is the client API for Lbwechat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LbwechatClient interface {
	// @cat:
	// @name:
	// @desc: 目前只做接口层面的定义，具体入参和出参待改造
	// @error:
	HandleWxGzhAuth(ctx context.Context, in *HandleWxGzhAuthReq, opts ...grpc.CallOption) (*HandleWxGzhAuthRsp, error)
	// @cat:
	// @name:
	// @desc: 目前只做接口层面的定义，具体入参和出参待改造
	// @error:
	HandleWxGzhMsg(ctx context.Context, in *HandleWxGzhMsgReq, opts ...grpc.CallOption) (*HandleWxGzhMsgRsp, error)
}

type lbwechatClient struct {
	cc grpc.ClientConnInterface
}

func NewLbwechatClient(cc grpc.ClientConnInterface) LbwechatClient {
	return &lbwechatClient{cc}
}

func (c *lbwechatClient) HandleWxGzhAuth(ctx context.Context, in *HandleWxGzhAuthReq, opts ...grpc.CallOption) (*HandleWxGzhAuthRsp, error) {
	out := new(HandleWxGzhAuthRsp)
	err := c.cc.Invoke(ctx, "/lbwechat.lbwechat/HandleWxGzhAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbwechatClient) HandleWxGzhMsg(ctx context.Context, in *HandleWxGzhMsgReq, opts ...grpc.CallOption) (*HandleWxGzhMsgRsp, error) {
	out := new(HandleWxGzhMsgRsp)
	err := c.cc.Invoke(ctx, "/lbwechat.lbwechat/HandleWxGzhMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LbwechatServer is the server API for Lbwechat service.
// All implementations must embed UnimplementedLbwechatServer
// for forward compatibility
type LbwechatServer interface {
	// @cat:
	// @name:
	// @desc: 目前只做接口层面的定义，具体入参和出参待改造
	// @error:
	HandleWxGzhAuth(context.Context, *HandleWxGzhAuthReq) (*HandleWxGzhAuthRsp, error)
	// @cat:
	// @name:
	// @desc: 目前只做接口层面的定义，具体入参和出参待改造
	// @error:
	HandleWxGzhMsg(context.Context, *HandleWxGzhMsgReq) (*HandleWxGzhMsgRsp, error)
	mustEmbedUnimplementedLbwechatServer()
}

// UnimplementedLbwechatServer must be embedded to have forward compatible implementations.
type UnimplementedLbwechatServer struct {
}

func (UnimplementedLbwechatServer) HandleWxGzhAuth(context.Context, *HandleWxGzhAuthReq) (*HandleWxGzhAuthRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleWxGzhAuth not implemented")
}
func (UnimplementedLbwechatServer) HandleWxGzhMsg(context.Context, *HandleWxGzhMsgReq) (*HandleWxGzhMsgRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleWxGzhMsg not implemented")
}
func (UnimplementedLbwechatServer) mustEmbedUnimplementedLbwechatServer() {}

// UnsafeLbwechatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LbwechatServer will
// result in compilation errors.
type UnsafeLbwechatServer interface {
	mustEmbedUnimplementedLbwechatServer()
}

func RegisterLbwechatServer(s grpc.ServiceRegistrar, srv LbwechatServer) {
	s.RegisterService(&Lbwechat_ServiceDesc, srv)
}

func _Lbwechat_HandleWxGzhAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleWxGzhAuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbwechatServer).HandleWxGzhAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbwechat.lbwechat/HandleWxGzhAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbwechatServer).HandleWxGzhAuth(ctx, req.(*HandleWxGzhAuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbwechat_HandleWxGzhMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleWxGzhMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbwechatServer).HandleWxGzhMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbwechat.lbwechat/HandleWxGzhMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbwechatServer).HandleWxGzhMsg(ctx, req.(*HandleWxGzhMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Lbwechat_ServiceDesc is the grpc.ServiceDesc for Lbwechat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lbwechat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lbwechat.lbwechat",
	HandlerType: (*LbwechatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleWxGzhAuth",
			Handler:    _Lbwechat_HandleWxGzhAuth_Handler,
		},
		{
			MethodName: "HandleWxGzhMsg",
			Handler:    _Lbwechat_HandleWxGzhMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbwechat.proto",
}
