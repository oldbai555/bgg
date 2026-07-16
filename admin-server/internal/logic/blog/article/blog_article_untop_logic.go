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

type BlogArticleUntopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleUntopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleUntopLogic {
	return &BlogArticleUntopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleUntop 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticleuntoplogic.go。
func (l *BlogArticleUntopLogic) BlogArticleUntop(req *types.BlogArticleUntopReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleUntop(l.ctx, &contentclient.BlogArticleUntopRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("取消置顶失败", err)
	}
	return &types.Response{Code: 0, Message: "取消置顶成功"}, nil
}
