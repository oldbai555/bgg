// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tag

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *BlogTagDeleteLogic) BlogTagDelete(req *types.BlogTagDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "标签ID不能为空")
	}

	// TODO: 如需限制有文章关联的标签删除，可以在后续通过文章标签关联表进行检查

	if err := blogrepo.NewBlogTagRepository(l.svcCtx.Repository).Delete(l.ctx, req.Id); err != nil {
		return err
	}

	return nil
}
