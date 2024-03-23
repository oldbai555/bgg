// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: lbblog.proto

package lbblog

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

// LbblogClient is the client API for Lbblog service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LbblogClient interface {
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddArticle(ctx context.Context, in *AddArticleReq, opts ...grpc.CallOption) (*AddArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelArticleList(ctx context.Context, in *DelArticleListReq, opts ...grpc.CallOption) (*DelArticleListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateArticle(ctx context.Context, in *UpdateArticleReq, opts ...grpc.CallOption) (*UpdateArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetArticle(ctx context.Context, in *GetArticleReq, opts ...grpc.CallOption) (*GetArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetArticleList(ctx context.Context, in *GetArticleListReq, opts ...grpc.CallOption) (*GetArticleListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddCategory(ctx context.Context, in *AddCategoryReq, opts ...grpc.CallOption) (*AddCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelCategoryList(ctx context.Context, in *DelCategoryListReq, opts ...grpc.CallOption) (*DelCategoryListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateCategory(ctx context.Context, in *UpdateCategoryReq, opts ...grpc.CallOption) (*UpdateCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCategory(ctx context.Context, in *GetCategoryReq, opts ...grpc.CallOption) (*GetCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCategoryList(ctx context.Context, in *GetCategoryListReq, opts ...grpc.CallOption) (*GetCategoryListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddComment(ctx context.Context, in *AddCommentReq, opts ...grpc.CallOption) (*AddCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelCommentList(ctx context.Context, in *DelCommentListReq, opts ...grpc.CallOption) (*DelCommentListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateComment(ctx context.Context, in *UpdateCommentReq, opts ...grpc.CallOption) (*UpdateCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetComment(ctx context.Context, in *GetCommentReq, opts ...grpc.CallOption) (*GetCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCommentList(ctx context.Context, in *GetCommentListReq, opts ...grpc.CallOption) (*GetCommentListRsp, error)
}

type lbblogClient struct {
	cc grpc.ClientConnInterface
}

func NewLbblogClient(cc grpc.ClientConnInterface) LbblogClient {
	return &lbblogClient{cc}
}

func (c *lbblogClient) AddArticle(ctx context.Context, in *AddArticleReq, opts ...grpc.CallOption) (*AddArticleRsp, error) {
	out := new(AddArticleRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/AddArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) DelArticleList(ctx context.Context, in *DelArticleListReq, opts ...grpc.CallOption) (*DelArticleListRsp, error) {
	out := new(DelArticleListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/DelArticleList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) UpdateArticle(ctx context.Context, in *UpdateArticleReq, opts ...grpc.CallOption) (*UpdateArticleRsp, error) {
	out := new(UpdateArticleRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/UpdateArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetArticle(ctx context.Context, in *GetArticleReq, opts ...grpc.CallOption) (*GetArticleRsp, error) {
	out := new(GetArticleRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetArticleList(ctx context.Context, in *GetArticleListReq, opts ...grpc.CallOption) (*GetArticleListRsp, error) {
	out := new(GetArticleListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetArticleList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) AddCategory(ctx context.Context, in *AddCategoryReq, opts ...grpc.CallOption) (*AddCategoryRsp, error) {
	out := new(AddCategoryRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/AddCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) DelCategoryList(ctx context.Context, in *DelCategoryListReq, opts ...grpc.CallOption) (*DelCategoryListRsp, error) {
	out := new(DelCategoryListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/DelCategoryList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) UpdateCategory(ctx context.Context, in *UpdateCategoryReq, opts ...grpc.CallOption) (*UpdateCategoryRsp, error) {
	out := new(UpdateCategoryRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/UpdateCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetCategory(ctx context.Context, in *GetCategoryReq, opts ...grpc.CallOption) (*GetCategoryRsp, error) {
	out := new(GetCategoryRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetCategoryList(ctx context.Context, in *GetCategoryListReq, opts ...grpc.CallOption) (*GetCategoryListRsp, error) {
	out := new(GetCategoryListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetCategoryList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) AddComment(ctx context.Context, in *AddCommentReq, opts ...grpc.CallOption) (*AddCommentRsp, error) {
	out := new(AddCommentRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/AddComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) DelCommentList(ctx context.Context, in *DelCommentListReq, opts ...grpc.CallOption) (*DelCommentListRsp, error) {
	out := new(DelCommentListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/DelCommentList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) UpdateComment(ctx context.Context, in *UpdateCommentReq, opts ...grpc.CallOption) (*UpdateCommentRsp, error) {
	out := new(UpdateCommentRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/UpdateComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetComment(ctx context.Context, in *GetCommentReq, opts ...grpc.CallOption) (*GetCommentRsp, error) {
	out := new(GetCommentRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lbblogClient) GetCommentList(ctx context.Context, in *GetCommentListReq, opts ...grpc.CallOption) (*GetCommentListRsp, error) {
	out := new(GetCommentListRsp)
	err := c.cc.Invoke(ctx, "/lbblog.lbblog/GetCommentList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LbblogServer is the server API for Lbblog service.
// All implementations must embed UnimplementedLbblogServer
// for forward compatibility
type LbblogServer interface {
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddArticle(context.Context, *AddArticleReq) (*AddArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelArticleList(context.Context, *DelArticleListReq) (*DelArticleListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateArticle(context.Context, *UpdateArticleReq) (*UpdateArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetArticle(context.Context, *GetArticleReq) (*GetArticleRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetArticleList(context.Context, *GetArticleListReq) (*GetArticleListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddCategory(context.Context, *AddCategoryReq) (*AddCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelCategoryList(context.Context, *DelCategoryListReq) (*DelCategoryListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateCategory(context.Context, *UpdateCategoryReq) (*UpdateCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCategory(context.Context, *GetCategoryReq) (*GetCategoryRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCategoryList(context.Context, *GetCategoryListReq) (*GetCategoryListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	AddComment(context.Context, *AddCommentReq) (*AddCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	DelCommentList(context.Context, *DelCommentListReq) (*DelCommentListRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	UpdateComment(context.Context, *UpdateCommentReq) (*UpdateCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetComment(context.Context, *GetCommentReq) (*GetCommentRsp, error)
	// @cat:
	// @name:
	// @desc:
	// @error:
	GetCommentList(context.Context, *GetCommentListReq) (*GetCommentListRsp, error)
	mustEmbedUnimplementedLbblogServer()
}

// UnimplementedLbblogServer must be embedded to have forward compatible implementations.
type UnimplementedLbblogServer struct {
}

func (UnimplementedLbblogServer) AddArticle(context.Context, *AddArticleReq) (*AddArticleRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddArticle not implemented")
}
func (UnimplementedLbblogServer) DelArticleList(context.Context, *DelArticleListReq) (*DelArticleListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelArticleList not implemented")
}
func (UnimplementedLbblogServer) UpdateArticle(context.Context, *UpdateArticleReq) (*UpdateArticleRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArticle not implemented")
}
func (UnimplementedLbblogServer) GetArticle(context.Context, *GetArticleReq) (*GetArticleRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticle not implemented")
}
func (UnimplementedLbblogServer) GetArticleList(context.Context, *GetArticleListReq) (*GetArticleListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleList not implemented")
}
func (UnimplementedLbblogServer) AddCategory(context.Context, *AddCategoryReq) (*AddCategoryRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCategory not implemented")
}
func (UnimplementedLbblogServer) DelCategoryList(context.Context, *DelCategoryListReq) (*DelCategoryListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelCategoryList not implemented")
}
func (UnimplementedLbblogServer) UpdateCategory(context.Context, *UpdateCategoryReq) (*UpdateCategoryRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCategory not implemented")
}
func (UnimplementedLbblogServer) GetCategory(context.Context, *GetCategoryReq) (*GetCategoryRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategory not implemented")
}
func (UnimplementedLbblogServer) GetCategoryList(context.Context, *GetCategoryListReq) (*GetCategoryListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategoryList not implemented")
}
func (UnimplementedLbblogServer) AddComment(context.Context, *AddCommentReq) (*AddCommentRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddComment not implemented")
}
func (UnimplementedLbblogServer) DelCommentList(context.Context, *DelCommentListReq) (*DelCommentListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelCommentList not implemented")
}
func (UnimplementedLbblogServer) UpdateComment(context.Context, *UpdateCommentReq) (*UpdateCommentRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateComment not implemented")
}
func (UnimplementedLbblogServer) GetComment(context.Context, *GetCommentReq) (*GetCommentRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComment not implemented")
}
func (UnimplementedLbblogServer) GetCommentList(context.Context, *GetCommentListReq) (*GetCommentListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentList not implemented")
}
func (UnimplementedLbblogServer) mustEmbedUnimplementedLbblogServer() {}

// UnsafeLbblogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LbblogServer will
// result in compilation errors.
type UnsafeLbblogServer interface {
	mustEmbedUnimplementedLbblogServer()
}

func RegisterLbblogServer(s grpc.ServiceRegistrar, srv LbblogServer) {
	s.RegisterService(&Lbblog_ServiceDesc, srv)
}

func _Lbblog_AddArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddArticleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).AddArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/AddArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).AddArticle(ctx, req.(*AddArticleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_DelArticleList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelArticleListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).DelArticleList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/DelArticleList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).DelArticleList(ctx, req.(*DelArticleListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_UpdateArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateArticleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).UpdateArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/UpdateArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).UpdateArticle(ctx, req.(*UpdateArticleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetArticle(ctx, req.(*GetArticleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetArticleList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetArticleList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetArticleList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetArticleList(ctx, req.(*GetArticleListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_AddCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCategoryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).AddCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/AddCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).AddCategory(ctx, req.(*AddCategoryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_DelCategoryList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelCategoryListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).DelCategoryList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/DelCategoryList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).DelCategoryList(ctx, req.(*DelCategoryListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_UpdateCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCategoryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).UpdateCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/UpdateCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).UpdateCategory(ctx, req.(*UpdateCategoryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetCategory(ctx, req.(*GetCategoryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetCategoryList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetCategoryList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetCategoryList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetCategoryList(ctx, req.(*GetCategoryListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_AddComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).AddComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/AddComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).AddComment(ctx, req.(*AddCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_DelCommentList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelCommentListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).DelCommentList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/DelCommentList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).DelCommentList(ctx, req.(*DelCommentListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_UpdateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).UpdateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/UpdateComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).UpdateComment(ctx, req.(*UpdateCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetComment(ctx, req.(*GetCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lbblog_GetCommentList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbblogServer).GetCommentList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbblog.lbblog/GetCommentList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbblogServer).GetCommentList(ctx, req.(*GetCommentListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Lbblog_ServiceDesc is the grpc.ServiceDesc for Lbblog service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lbblog_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lbblog.lbblog",
	HandlerType: (*LbblogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddArticle",
			Handler:    _Lbblog_AddArticle_Handler,
		},
		{
			MethodName: "DelArticleList",
			Handler:    _Lbblog_DelArticleList_Handler,
		},
		{
			MethodName: "UpdateArticle",
			Handler:    _Lbblog_UpdateArticle_Handler,
		},
		{
			MethodName: "GetArticle",
			Handler:    _Lbblog_GetArticle_Handler,
		},
		{
			MethodName: "GetArticleList",
			Handler:    _Lbblog_GetArticleList_Handler,
		},
		{
			MethodName: "AddCategory",
			Handler:    _Lbblog_AddCategory_Handler,
		},
		{
			MethodName: "DelCategoryList",
			Handler:    _Lbblog_DelCategoryList_Handler,
		},
		{
			MethodName: "UpdateCategory",
			Handler:    _Lbblog_UpdateCategory_Handler,
		},
		{
			MethodName: "GetCategory",
			Handler:    _Lbblog_GetCategory_Handler,
		},
		{
			MethodName: "GetCategoryList",
			Handler:    _Lbblog_GetCategoryList_Handler,
		},
		{
			MethodName: "AddComment",
			Handler:    _Lbblog_AddComment_Handler,
		},
		{
			MethodName: "DelCommentList",
			Handler:    _Lbblog_DelCommentList_Handler,
		},
		{
			MethodName: "UpdateComment",
			Handler:    _Lbblog_UpdateComment_Handler,
		},
		{
			MethodName: "GetComment",
			Handler:    _Lbblog_GetComment_Handler,
		},
		{
			MethodName: "GetCommentList",
			Handler:    _Lbblog_GetCommentList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbblog.proto",
}
