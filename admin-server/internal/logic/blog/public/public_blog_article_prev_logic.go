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

type PublicBlogArticlePrevLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticlePrevLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticlePrevLogic {
	return &PublicBlogArticlePrevLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogArticlePrev 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogarticleprevlogic.go。
func (l *PublicBlogArticlePrevLogic) PublicBlogArticlePrev(req *types.PublicBlogArticlePrevReq) (resp *types.PublicBlogArticlePrevResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogArticlePrev(l.ctx, &contentclient.PublicBlogArticlePrevNextRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询上一篇文章失败", err)
	}
	return &types.PublicBlogArticlePrevResp{Id: rpcResp.Id, Title: rpcResp.Title, PublishTime: rpcResp.PublishTime}, nil
}
