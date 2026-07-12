// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUpdateLogic {
	return &BlogArticleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleUpdate 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticleupdatelogic.go。
func (l *BlogArticleUpdateLogic) BlogArticleUpdate(req *types.BlogArticleUpdateReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleUpdate(l.ctx, &contentclient.BlogArticleUpdateRequest{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		TagIds:  req.TagIds,
		Cover:   req.Cover,
		Summary: req.Summary,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("更新文章失败", err)
	}
	return &types.Response{Code: int(errs.CodeOK), Message: "更新成功"}, nil
}
