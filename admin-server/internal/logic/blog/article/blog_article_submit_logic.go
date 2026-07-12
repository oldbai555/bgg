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

type BlogArticleSubmitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleSubmitLogic {
	return &BlogArticleSubmitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleSubmit 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticlesubmitlogic.go。
func (l *BlogArticleSubmitLogic) BlogArticleSubmit(req *types.BlogArticleSubmitReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticleSubmit(l.ctx, &contentclient.BlogArticleSubmitRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("提交审核失败", err)
	}
	return &types.Response{Code: int(errs.CodeOK), Message: "已提交审核"}, nil
}
