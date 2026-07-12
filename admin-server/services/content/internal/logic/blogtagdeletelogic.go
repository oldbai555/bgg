package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagDeleteLogic {
	return &BlogTagDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogTagDelete 迁移自 internal/logic/blog/tag/blog_tag_delete_logic.go。
func (l *BlogTagDeleteLogic) BlogTagDelete(in *content.BlogTagDeleteRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "标签ID不能为空"))
	}
	if err := l.svcCtx.Tag.Delete(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &content.Empty{}, nil
}
