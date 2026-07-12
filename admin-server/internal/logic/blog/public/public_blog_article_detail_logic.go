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

type PublicBlogArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleDetailLogic {
	return &PublicBlogArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogArticleDetail 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogarticledetaillogic.go。
func (l *PublicBlogArticleDetailLogic) PublicBlogArticleDetail(req *types.PublicBlogArticleDetailReq) (resp *types.PublicBlogArticleDetailResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogArticleDetail(l.ctx, &contentclient.PublicBlogArticleDetailRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询文章详情失败", err)
	}

	tagItems := make([]types.BlogTagItem, 0, len(rpcResp.Tags))
	for _, t := range rpcResp.Tags {
		tagItems = append(tagItems, types.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &types.PublicBlogArticleDetailResp{
		Id:          rpcResp.Id,
		Title:       rpcResp.Title,
		Content:     rpcResp.Content,
		Cover:       rpcResp.Cover,
		AuthorName:  rpcResp.AuthorName,
		PublishTime: rpcResp.PublishTime,
		Tags:        tagItems,
	}, nil
}
