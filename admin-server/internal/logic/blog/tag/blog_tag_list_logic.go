// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tag

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagListLogic {
	return &BlogTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogTagList 薄胶水：解析 HTTP 请求 -> 调 ContentRPC -> 映射响应，实际业务逻辑已经搬进
// services/content/internal/logic/blogtaglistlogic.go。
func (l *BlogTagListLogic) BlogTagList(req *types.BlogTagListReq) (resp *types.BlogTagListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.BlogTagList(l.ctx, &contentclient.BlogTagListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Status:   req.Status,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询标签列表失败", err)
	}

	items := make([]types.BlogTagItem, 0, len(rpcResp.List))
	for _, t := range rpcResp.List {
		items = append(items, types.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &types.BlogTagListResp{Total: rpcResp.Total, List: items}, nil
}
