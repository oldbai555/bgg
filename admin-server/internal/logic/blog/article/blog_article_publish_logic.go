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

type BlogArticlePublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticlePublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticlePublishLogic {
	return &BlogArticlePublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticlePublish 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticlepublishlogic.go。
func (l *BlogArticlePublishLogic) BlogArticlePublish(req *types.BlogArticlePublishReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogArticlePublish(l.ctx, &contentclient.BlogArticlePublishRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("上架失败", err)
	}
	return &types.Response{Code: int(errs.CodeOK), Message: "上架成功"}, nil
}
