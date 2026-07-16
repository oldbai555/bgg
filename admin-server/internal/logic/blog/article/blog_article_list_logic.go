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

type BlogArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleListLogic {
	return &BlogArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticlelistlogic.go。
func (l *BlogArticleListLogic) BlogArticleList(req *types.BlogArticleListReq) (resp *types.BlogArticleListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.BlogArticleList(l.ctx, &contentclient.BlogArticleListRequest{
		Page:        req.Page,
		PageSize:    req.PageSize,
		Title:       req.Title,
		Status:      req.Status,
		AuditStatus: req.AuditStatus,
		TagId:       req.TagId,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询文章列表失败", err)
	}

	items := make([]types.BlogArticleItem, 0, len(rpcResp.List))
	for _, a := range rpcResp.List {
		items = append(items, types.BlogArticleItem{
			Id:          a.Id,
			Title:       a.Title,
			Status:      a.Status,
			AuditStatus: a.AuditStatus,
			Cover:       a.Cover,
			AuthorId:    a.AuthorId,
			AuthorName:  a.AuthorName,
			TagIds:      a.TagIds,
			TagNames:    a.TagNames,
			PublishTime: a.PublishTime,
			IsTop:       a.IsTop,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		})
	}

	return &types.BlogArticleListResp{Total: rpcResp.Total, List: items}, nil
}
