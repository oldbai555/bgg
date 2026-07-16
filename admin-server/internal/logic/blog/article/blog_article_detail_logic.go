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

type BlogArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleDetailLogic {
	return &BlogArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleDetail 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticledetaillogic.go。
func (l *BlogArticleDetailLogic) BlogArticleDetail(req *types.BlogArticleDetailReq) (resp *types.BlogArticleDetailResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.BlogArticleDetail(l.ctx, &contentclient.BlogArticleDetailRequest{Id: req.Id})
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

	return &types.BlogArticleDetailResp{
		Id:          rpcResp.Id,
		Title:       rpcResp.Title,
		Content:     rpcResp.Content,
		Status:      rpcResp.Status,
		AuditStatus: rpcResp.AuditStatus,
		Cover:       rpcResp.Cover,
		AuthorId:    rpcResp.AuthorId,
		AuthorName:  rpcResp.AuthorName,
		PublishTime: rpcResp.PublishTime,
		Summary:     rpcResp.Summary,
		Tags:        tagItems,
		CreatedAt:   rpcResp.CreatedAt,
		UpdatedAt:   rpcResp.UpdatedAt,
	}, nil
}
