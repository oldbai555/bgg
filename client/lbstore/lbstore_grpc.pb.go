// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: lbstore.proto

package lbstore

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

// LbstoreClient is the client API for Lbstore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LbstoreClient interface {
	// @cat:
	// @name:
	// @desc:
	// @error:
	Upload(ctx context.Context, in *UploadReq, opts ...grpc.CallOption) (*UploadRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	RefreshFileSignedUrl(ctx context.Context, in *RefreshFileSignedUrlReq, opts ...grpc.CallOption) (*RefreshFileSignedUrlRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetSignature(ctx context.Context, in *GetSignatureReq, opts ...grpc.CallOption) (*GetSignatureRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	ReportUploadFile(ctx context.Context, in *ReportUploadFileReq, opts ...grpc.CallOption) (*ReportUploadFileRsp, error)
}

type lbstoreClient struct {
	cc grpc.ClientConnInterface
}

func NewLbstoreClient(cc grpc.ClientConnInterface) LbstoreClient {
	return &lbstoreClient{cc}
}

func (c *lbstoreClient) Upload(ctx context.Context, in *UploadReq, opts ...grpc.CallOption) (*UploadRsp, error) {
	out := new(UploadRsp)
	err := c.cc.Invoke(ctx, "/lbstore.lbstore/Upload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbstoreClient) GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListRsp, error) {
	out := new(GetFileListRsp)
	err := c.cc.Invoke(ctx, "/lbstore.lbstore/GetFileList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbstoreClient) RefreshFileSignedUrl(ctx context.Context, in *RefreshFileSignedUrlReq, opts ...grpc.CallOption) (*RefreshFileSignedUrlRsp, error) {
	out := new(RefreshFileSignedUrlRsp)
	err := c.cc.Invoke(ctx, "/lbstore.lbstore/RefreshFileSignedUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbstoreClient) GetSignature(ctx context.Context, in *GetSignatureReq, opts ...grpc.CallOption) (*GetSignatureRsp, error) {
	out := new(GetSignatureRsp)
	err := c.cc.Invoke(ctx, "/lbstore.lbstore/GetSignature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbstoreClient) ReportUploadFile(ctx context.Context, in *ReportUploadFileReq, opts ...grpc.CallOption) (*ReportUploadFileRsp, error) {
	out := new(ReportUploadFileRsp)
	err := c.cc.Invoke(ctx, "/lbstore.lbstore/ReportUploadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LbstoreServer is the server API for Lbstore service.
// All implementations must embed UnimplementedLbstoreServer
// for forward compatibility
type LbstoreServer interface {
	// @cat:
	// @name:
	// @desc:
	// @error:
	Upload(context.Context, *UploadReq) (*UploadRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetFileList(context.Context, *GetFileListReq) (*GetFileListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	RefreshFileSignedUrl(context.Context, *RefreshFileSignedUrlReq) (*RefreshFileSignedUrlRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetSignature(context.Context, *GetSignatureReq) (*GetSignatureRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	ReportUploadFile(context.Context, *ReportUploadFileReq) (*ReportUploadFileRsp, error)
	mustEmbedUnimplementedLbstoreServer()
}

// UnimplementedLbstoreServer must be embedded to have forward compatible implementations.
type UnimplementedLbstoreServer struct {
}

func (UnimplementedLbstoreServer) Upload(context.Context, *UploadReq) (*UploadRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedLbstoreServer) GetFileList(context.Context, *GetFileListReq) (*GetFileListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileList not implemented")
}
func (UnimplementedLbstoreServer) RefreshFileSignedUrl(context.Context, *RefreshFileSignedUrlReq) (*RefreshFileSignedUrlRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshFileSignedUrl not implemented")
}
func (UnimplementedLbstoreServer) GetSignature(context.Context, *GetSignatureReq) (*GetSignatureRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSignature not implemented")
}
func (UnimplementedLbstoreServer) ReportUploadFile(context.Context, *ReportUploadFileReq) (*ReportUploadFileRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportUploadFile not implemented")
}
func (UnimplementedLbstoreServer) mustEmbedUnimplementedLbstoreServer() {}

// UnsafeLbstoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LbstoreServer will
// result in compilation errors.
type UnsafeLbstoreServer interface {
	mustEmbedUnimplementedLbstoreServer()
}

func RegisterLbstoreServer(s grpc.ServiceRegistrar, srv LbstoreServer) {
	s.RegisterService(&Lbstore_ServiceDesc, srv)
}

func _Lbstore_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbstoreServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbstore.lbstore/Upload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbstoreServer).Upload(ctx, req.(*UploadReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbstore_GetFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbstoreServer).GetFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbstore.lbstore/GetFileList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbstoreServer).GetFileList(ctx, req.(*GetFileListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbstore_RefreshFileSignedUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshFileSignedUrlReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbstoreServer).RefreshFileSignedUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbstore.lbstore/RefreshFileSignedUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbstoreServer).RefreshFileSignedUrl(ctx, req.(*RefreshFileSignedUrlReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbstore_GetSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSignatureReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbstoreServer).GetSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbstore.lbstore/GetSignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbstoreServer).GetSignature(ctx, req.(*GetSignatureReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbstore_ReportUploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportUploadFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbstoreServer).ReportUploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbstore.lbstore/ReportUploadFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbstoreServer).ReportUploadFile(ctx, req.(*ReportUploadFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Lbstore_ServiceDesc is the grpc.ServiceDesc for Lbstore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lbstore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lbstore.lbstore",
	HandlerType: (*LbstoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upload",
			Handler:    _Lbstore_Upload_Handler,
		},
		{
			MethodName: "GetFileList",
			Handler:    _Lbstore_GetFileList_Handler,
		},
		{
			MethodName: "RefreshFileSignedUrl",
			Handler:    _Lbstore_RefreshFileSignedUrl_Handler,
		},
		{
			MethodName: "GetSignature",
			Handler:    _Lbstore_GetSignature_Handler,
		},
		{
			MethodName: "ReportUploadFile",
			Handler:    _Lbstore_ReportUploadFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbstore.proto",
}
