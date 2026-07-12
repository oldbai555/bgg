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

type BlogArticleTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleTopLogic {
	return &BlogArticleTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleTop 薄胶水，置顶数量上限的静态配置已经下沉到 content-rpc，实际业务逻辑
// 搬进 services/content/internal/logic/blogarticletoplogic.go。
func (l *BlogArticleTopLogic) BlogArticleTop(req *types.BlogArticleTopReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleTop(l.ctx, &contentclient.BlogArticleTopRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("置顶失败", err)
	}
	return &types.Response{Code: 0, Message: "置顶成功"}, nil
}
