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

type BlogTagCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagCreateLogic {
	return &BlogTagCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogTagCreate 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/blogtagcreatelogic.go。
func (l *BlogTagCreateLogic) BlogTagCreate(req *types.BlogTagCreateReq) error {
	_, err := l.svcCtx.ContentRPC.BlogTagCreate(l.ctx, &contentclient.BlogTagCreateRequest{
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("创建标签失败", err)
	}
	return nil
}
