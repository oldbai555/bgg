// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: lbsingle.proto

package client

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

// LbsingleClient is the client API for Lbsingle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LbsingleClient interface {
	// @cat:
	// @name:
	// @desc: 新增文件
	// @error:
	AddFile(ctx context.Context, in *AddFileReq, opts ...grpc.CallOption) (*AddFileRsp, error)
	// @cat:
	// @name:
	// @desc: 删除文件
	// @error:
	DelFileList(ctx context.Context, in *DelFileListReq, opts ...grpc.CallOption) (*DelFileListRsp, error)
	// @cat:
	// @name:
	// @desc: 更新文件
	// @error:
	UpdateFile(ctx context.Context, in *UpdateFileReq, opts ...grpc.CallOption) (*UpdateFileRsp, error)
	// @cat:
	// @name:
	// @desc: 获取单个文件
	// @error:
	GetFile(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (*GetFileRsp, error)
	// @cat:
	// @name:
	// @desc: 获取文件列表
	// @error:
	GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListRsp, error)
	// @desc: 登录
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRsp, error)
	// @desc: 登出
	Logout(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*LogoutRsp, error)
	// @desc: 获取登录用户的信息
	GetLoginUser(ctx context.Context, in *GetLoginUserReq, opts ...grpc.CallOption) (*GetLoginUserRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error: 更新登陆的用户信息
	UpdateLoginUserInfo(ctx context.Context, in *UpdateLoginUserInfoReq, opts ...grpc.CallOption) (*UpdateLoginUserInfoRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	ResetPassword(ctx context.Context, in *ResetPasswordReq, opts ...grpc.CallOption) (*ResetPasswordRsp, error)
}

type lbsingleClient struct {
	cc grpc.ClientConnInterface
}

func NewLbsingleClient(cc grpc.ClientConnInterface) LbsingleClient {
	return &lbsingleClient{cc}
}

func (c *lbsingleClient) AddFile(ctx context.Context, in *AddFileReq, opts ...grpc.CallOption) (*AddFileRsp, error) {
	out := new(AddFileRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/AddFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) DelFileList(ctx context.Context, in *DelFileListReq, opts ...grpc.CallOption) (*DelFileListRsp, error) {
	out := new(DelFileListRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/DelFileList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) UpdateFile(ctx context.Context, in *UpdateFileReq, opts ...grpc.CallOption) (*UpdateFileRsp, error) {
	out := new(UpdateFileRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/UpdateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) GetFile(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (*GetFileRsp, error) {
	out := new(GetFileRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/GetFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListRsp, error) {
	out := new(GetFileListRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/GetFileList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRsp, error) {
	out := new(LoginRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) Logout(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*LogoutRsp, error) {
	out := new(LogoutRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) GetLoginUser(ctx context.Context, in *GetLoginUserReq, opts ...grpc.CallOption) (*GetLoginUserRsp, error) {
	out := new(GetLoginUserRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/GetLoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) UpdateLoginUserInfo(ctx context.Context, in *UpdateLoginUserInfoReq, opts ...grpc.CallOption) (*UpdateLoginUserInfoRsp, error) {
	out := new(UpdateLoginUserInfoRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/UpdateLoginUserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbsingleClient) ResetPassword(ctx context.Context, in *ResetPasswordReq, opts ...grpc.CallOption) (*ResetPasswordRsp, error) {
	out := new(ResetPasswordRsp)
	err := c.cc.Invoke(ctx, "/lbsingle.lbsingle/ResetPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LbsingleServer is the server API for Lbsingle service.
// All implementations must embed UnimplementedLbsingleServer
// for forward compatibility
type LbsingleServer interface {
	// @cat:
	// @name:
	// @desc: 新增文件
	// @error:
	AddFile(context.Context, *AddFileReq) (*AddFileRsp, error)
	// @cat:
	// @name:
	// @desc: 删除文件
	// @error:
	DelFileList(context.Context, *DelFileListReq) (*DelFileListRsp, error)
	// @cat:
	// @name:
	// @desc: 更新文件
	// @error:
	UpdateFile(context.Context, *UpdateFileReq) (*UpdateFileRsp, error)
	// @cat:
	// @name:
	// @desc: 获取单个文件
	// @error:
	GetFile(context.Context, *GetFileReq) (*GetFileRsp, error)
	// @cat:
	// @name:
	// @desc: 获取文件列表
	// @error:
	GetFileList(context.Context, *GetFileListReq) (*GetFileListRsp, error)
	// @desc: 登录
	Login(context.Context, *LoginReq) (*LoginRsp, error)
	// @desc: 登出
	Logout(context.Context, *LogoutReq) (*LogoutRsp, error)
	// @desc: 获取登录用户的信息
	GetLoginUser(context.Context, *GetLoginUserReq) (*GetLoginUserRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error: 更新登陆的用户信息
	UpdateLoginUserInfo(context.Context, *UpdateLoginUserInfoReq) (*UpdateLoginUserInfoRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	ResetPassword(context.Context, *ResetPasswordReq) (*ResetPasswordRsp, error)
	mustEmbedUnimplementedLbsingleServer()
}

// UnimplementedLbsingleServer must be embedded to have forward compatible implementations.
type UnimplementedLbsingleServer struct {
}

func (UnimplementedLbsingleServer) AddFile(context.Context, *AddFileReq) (*AddFileRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}
func (UnimplementedLbsingleServer) DelFileList(context.Context, *DelFileListReq) (*DelFileListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelFileList not implemented")
}
func (UnimplementedLbsingleServer) UpdateFile(context.Context, *UpdateFileReq) (*UpdateFileRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFile not implemented")
}
func (UnimplementedLbsingleServer) GetFile(context.Context, *GetFileReq) (*GetFileRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedLbsingleServer) GetFileList(context.Context, *GetFileListReq) (*GetFileListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileList not implemented")
}
func (UnimplementedLbsingleServer) Login(context.Context, *LoginReq) (*LoginRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedLbsingleServer) Logout(context.Context, *LogoutReq) (*LogoutRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedLbsingleServer) GetLoginUser(context.Context, *GetLoginUserReq) (*GetLoginUserRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoginUser not implemented")
}
func (UnimplementedLbsingleServer) UpdateLoginUserInfo(context.Context, *UpdateLoginUserInfoReq) (*UpdateLoginUserInfoRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoginUserInfo not implemented")
}
func (UnimplementedLbsingleServer) ResetPassword(context.Context, *ResetPasswordReq) (*ResetPasswordRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPassword not implemented")
}
func (UnimplementedLbsingleServer) mustEmbedUnimplementedLbsingleServer() {}

// UnsafeLbsingleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LbsingleServer will
// result in compilation errors.
type UnsafeLbsingleServer interface {
	mustEmbedUnimplementedLbsingleServer()
}

func RegisterLbsingleServer(s grpc.ServiceRegistrar, srv LbsingleServer) {
	s.RegisterService(&Lbsingle_ServiceDesc, srv)
}

func _Lbsingle_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/AddFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).AddFile(ctx, req.(*AddFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_DelFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelFileListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).DelFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/DelFileList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).DelFileList(ctx, req.(*DelFileListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_UpdateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).UpdateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/UpdateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).UpdateFile(ctx, req.(*UpdateFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/GetFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).GetFile(ctx, req.(*GetFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_GetFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).GetFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/GetFileList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).GetFileList(ctx, req.(*GetFileListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).Logout(ctx, req.(*LogoutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_GetLoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLoginUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).GetLoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/GetLoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).GetLoginUser(ctx, req.(*GetLoginUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_UpdateLoginUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLoginUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).UpdateLoginUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/UpdateLoginUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).UpdateLoginUserInfo(ctx, req.(*UpdateLoginUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbsingle_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbsingleServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbsingle.lbsingle/ResetPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbsingleServer).ResetPassword(ctx, req.(*ResetPasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Lbsingle_ServiceDesc is the grpc.ServiceDesc for Lbsingle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lbsingle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lbsingle.lbsingle",
	HandlerType: (*LbsingleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddFile",
			Handler:    _Lbsingle_AddFile_Handler,
		},
		{
			MethodName: "DelFileList",
			Handler:    _Lbsingle_DelFileList_Handler,
		},
		{
			MethodName: "UpdateFile",
			Handler:    _Lbsingle_UpdateFile_Handler,
		},
		{
			MethodName: "GetFile",
			Handler:    _Lbsingle_GetFile_Handler,
		},
		{
			MethodName: "GetFileList",
			Handler:    _Lbsingle_GetFileList_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Lbsingle_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Lbsingle_Logout_Handler,
		},
		{
			MethodName: "GetLoginUser",
			Handler:    _Lbsingle_GetLoginUser_Handler,
		},
		{
			MethodName: "UpdateLoginUserInfo",
			Handler:    _Lbsingle_UpdateLoginUserInfo_Handler,
		},
		{
			MethodName: "ResetPassword",
			Handler:    _Lbsingle_ResetPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbsingle.proto",
}
