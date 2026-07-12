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

type BlogArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDeleteLogic {
	return &BlogArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleDelete 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticledeletelogic.go。
func (l *BlogArticleDeleteLogic) BlogArticleDelete(req *types.BlogArticleDeleteReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleDelete(l.ctx, &contentclient.BlogArticleDeleteRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("删除文章失败", err)
	}
	return &types.Response{Code: int(errs.CodeOK), Message: "删除成功"}, nil
}
