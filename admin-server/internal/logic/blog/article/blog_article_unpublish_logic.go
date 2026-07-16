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

type BlogArticleUnpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUnpublishLogic {
	return &BlogArticleUnpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleUnpublish 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticleunpublishlogic.go。
func (l *BlogArticleUnpublishLogic) BlogArticleUnpublish(req *types.BlogArticleUnpublishReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleUnpublish(l.ctx, &contentclient.BlogArticleUnpublishRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("下架失败", err)
	}
	return &types.Response{Code: int(errs.CodeOK), Message: "下架成功"}, nil
}
