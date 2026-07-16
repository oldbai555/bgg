// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleNextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleNextLogic {
	return &PublicBlogArticleNextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogArticleNext 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogarticlenextlogic.go。
func (l *PublicBlogArticleNextLogic) PublicBlogArticleNext(req *types.PublicBlogArticleNextReq) (resp *types.PublicBlogArticleNextResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogArticleNext(l.ctx, &contentclient.PublicBlogArticlePrevNextRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询下一篇文章失败", err)
	}
	return &types.PublicBlogArticleNextResp{Id: rpcResp.Id, Title: rpcResp.Title, PublishTime: rpcResp.PublishTime}, nil
}
