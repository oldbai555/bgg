package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbblog"
)

var lbblogServer LbblogServer

type LbblogServer struct {
	*lbblog.UnimplementedLbblogServer
}

func (a *LbblogServer) AddComment(ctx context.Context, req *lbblog.AddCommentReq) (*lbblog.AddCommentRsp, error) {
	var rsp lbblog.AddCommentRsp
	return &rsp, nil
}
func (a *LbblogServer) GetArticle(ctx context.Context, req *lbblog.GetArticleReq) (*lbblog.GetArticleRsp, error) {
	var rsp lbblog.GetArticleRsp
	return &rsp, nil
}
func (a *LbblogServer) DelArticle(ctx context.Context, req *lbblog.DelArticleReq) (*lbblog.DelArticleRsp, error) {
	var rsp lbblog.DelArticleRsp
	return &rsp, nil
}
func (a *LbblogServer) DelCategory(ctx context.Context, req *lbblog.DelCategoryReq) (*lbblog.DelCategoryRsp, error) {
	var rsp lbblog.DelCategoryRsp
	return &rsp, nil
}
func (a *LbblogServer) AddCategory(ctx context.Context, req *lbblog.AddCategoryReq) (*lbblog.AddCategoryRsp, error) {
	var rsp lbblog.AddCategoryRsp
	return &rsp, nil
}
func (a *LbblogServer) UpdateArticle(ctx context.Context, req *lbblog.UpdateArticleReq) (*lbblog.UpdateArticleRsp, error) {
	var rsp lbblog.UpdateArticleRsp
	return &rsp, nil
}
func (a *LbblogServer) GetComment(ctx context.Context, req *lbblog.GetCommentReq) (*lbblog.GetCommentRsp, error) {
	var rsp lbblog.GetCommentRsp
	return &rsp, nil
}
func (a *LbblogServer) UpdateComment(ctx context.Context, req *lbblog.UpdateCommentReq) (*lbblog.UpdateCommentRsp, error) {
	var rsp lbblog.UpdateCommentRsp
	return &rsp, nil
}
func (a *LbblogServer) GetArticleList(ctx context.Context, req *lbblog.GetArticleListReq) (*lbblog.GetArticleListRsp, error) {
	var rsp lbblog.GetArticleListRsp
	return &rsp, nil
}
func (a *LbblogServer) GetCategoryList(ctx context.Context, req *lbblog.GetCategoryListReq) (*lbblog.GetCategoryListRsp, error) {
	var rsp lbblog.GetCategoryListRsp
	return &rsp, nil
}
func (a *LbblogServer) GetCategory(ctx context.Context, req *lbblog.GetCategoryReq) (*lbblog.GetCategoryRsp, error) {
	var rsp lbblog.GetCategoryRsp
	return &rsp, nil
}
func (a *LbblogServer) UpdateCategory(ctx context.Context, req *lbblog.UpdateCategoryReq) (*lbblog.UpdateCategoryRsp, error) {
	var rsp lbblog.UpdateCategoryRsp
	return &rsp, nil
}
func (a *LbblogServer) AddArticle(ctx context.Context, req *lbblog.AddArticleReq) (*lbblog.AddArticleRsp, error) {
	var rsp lbblog.AddArticleRsp
	return &rsp, nil
}
func (a *LbblogServer) GetCommentList(ctx context.Context, req *lbblog.GetCommentListReq) (*lbblog.GetCommentListRsp, error) {
	var rsp lbblog.GetCommentListRsp
	return &rsp, nil
}
func (a *LbblogServer) DelComment(ctx context.Context, req *lbblog.DelCommentReq) (*lbblog.DelCommentRsp, error) {
	var rsp lbblog.DelCommentRsp
	return &rsp, nil
}
