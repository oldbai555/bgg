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

type BlogTagDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagDeleteLogic {
	return &BlogTagDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogTagDelete 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/blogtagdeletelogic.go。
func (l *BlogTagDeleteLogic) BlogTagDelete(req *types.BlogTagDeleteReq) error {
	_, err := l.svcCtx.ContentRPC.BlogTagDelete(l.ctx, &contentclient.BlogTagDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除标签失败", err)
	}
	return nil
}
