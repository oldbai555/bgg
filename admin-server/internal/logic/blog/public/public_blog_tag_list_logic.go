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

type PublicBlogTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogTagListLogic {
	return &PublicBlogTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogTagList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogtaglistlogic.go。
func (l *PublicBlogTagListLogic) PublicBlogTagList() (resp *types.PublicBlogTagListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogTagList(l.ctx, &contentclient.PublicBlogGlobalRequest{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询标签列表失败", err)
	}
	items := make([]types.BlogTagItem, 0, len(rpcResp.List))
	for _, tag := range rpcResp.List {
		items = append(items, types.BlogTagItem{
			Id:        tag.Id,
			Name:      tag.Name,
			Status:    tag.Status,
			Remark:    tag.Remark,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	return &types.PublicBlogTagListResp{List: items}, nil
}
