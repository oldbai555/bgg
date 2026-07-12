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

type BlogTagUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagUpdateLogic {
	return &BlogTagUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogTagUpdate 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/blogtagupdatelogic.go。
func (l *BlogTagUpdateLogic) BlogTagUpdate(req *types.BlogTagUpdateReq) error {
	_, err := l.svcCtx.ContentRPC.BlogTagUpdate(l.ctx, &contentclient.BlogTagUpdateRequest{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("更新标签失败", err)
	}
	return nil
}
