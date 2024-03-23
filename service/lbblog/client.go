// Code generated by baixctl genclient. DO NOT EDIT.
package lbblog

import (
	"context"
	"github.com/oldbai555/bgg/internal/_const"
	"github.com/oldbai555/bgg/internal/bhttp"
	"net/http"
)

const ClientName = "lbblog"

func AddArticle(ctx context.Context, req *AddArticleReq) (*AddArticleRsp, error) {
	var rsp AddArticleRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, AddArticleCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.AddArticle(ctx, req)
}

func DelArticleList(ctx context.Context, req *DelArticleListReq) (*DelArticleListRsp, error) {
	var rsp DelArticleListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, DelArticleListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.DelArticleList(ctx, req)
}

func UpdateArticle(ctx context.Context, req *UpdateArticleReq) (*UpdateArticleRsp, error) {
	var rsp UpdateArticleRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, UpdateArticleCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.UpdateArticle(ctx, req)
}

func GetArticle(ctx context.Context, req *GetArticleReq) (*GetArticleRsp, error) {
	var rsp GetArticleRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetArticleCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetArticle(ctx, req)
}

func GetArticleList(ctx context.Context, req *GetArticleListReq) (*GetArticleListRsp, error) {
	var rsp GetArticleListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetArticleListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetArticleList(ctx, req)
}

func AddCategory(ctx context.Context, req *AddCategoryReq) (*AddCategoryRsp, error) {
	var rsp AddCategoryRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, AddCategoryCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.AddCategory(ctx, req)
}

func DelCategoryList(ctx context.Context, req *DelCategoryListReq) (*DelCategoryListRsp, error) {
	var rsp DelCategoryListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, DelCategoryListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.DelCategoryList(ctx, req)
}

func UpdateCategory(ctx context.Context, req *UpdateCategoryReq) (*UpdateCategoryRsp, error) {
	var rsp UpdateCategoryRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, UpdateCategoryCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.UpdateCategory(ctx, req)
}

func GetCategory(ctx context.Context, req *GetCategoryReq) (*GetCategoryRsp, error) {
	var rsp GetCategoryRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetCategoryCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetCategory(ctx, req)
}

func GetCategoryList(ctx context.Context, req *GetCategoryListReq) (*GetCategoryListRsp, error) {
	var rsp GetCategoryListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetCategoryListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetCategoryList(ctx, req)
}

func AddComment(ctx context.Context, req *AddCommentReq) (*AddCommentRsp, error) {
	var rsp AddCommentRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, AddCommentCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.AddComment(ctx, req)
}

func DelCommentList(ctx context.Context, req *DelCommentListReq) (*DelCommentListRsp, error) {
	var rsp DelCommentListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, DelCommentListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.DelCommentList(ctx, req)
}

func UpdateComment(ctx context.Context, req *UpdateCommentReq) (*UpdateCommentRsp, error) {
	var rsp UpdateCommentRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, UpdateCommentCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.UpdateComment(ctx, req)
}

func GetComment(ctx context.Context, req *GetCommentReq) (*GetCommentRsp, error) {
	var rsp GetCommentRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetCommentCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetComment(ctx, req)
}

func GetCommentList(ctx context.Context, req *GetCommentListReq) (*GetCommentListRsp, error) {
	var rsp GetCommentListRsp
	if cliMgr.conn == nil {
		return &rsp, bhttp.DoRequest(ctx, ServerName, GetCommentListCMDPath, http.MethodPost, _const.PROTO_TYPE_PROTO3, req, &rsp)
	}
	return cliMgr.cli.GetCommentList(ctx, req)
}
