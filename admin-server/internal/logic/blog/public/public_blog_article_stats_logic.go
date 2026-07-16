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

type PublicBlogArticleStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleStatsLogic {
	return &PublicBlogArticleStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogArticleStats 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogarticlestatslogic.go。
func (l *PublicBlogArticleStatsLogic) PublicBlogArticleStats() (resp *types.PublicBlogArticleStatsResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogArticleStats(l.ctx, &contentclient.PublicBlogGlobalRequest{})
	if err != nil {
		return nil, errs.WrapGRPCError("统计文章总数失败", err)
	}
	return &types.PublicBlogArticleStatsResp{TotalArticles: rpcResp.TotalArticles}, nil
}
