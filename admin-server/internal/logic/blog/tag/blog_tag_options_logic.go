package tag

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

// BlogTagOptionsLogic 标签下拉选项逻辑
type BlogTagOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagOptionsLogic {
	return &BlogTagOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogTagOptionsLogic) BlogTagOptions(req *types.BlogTagOptionsReq) (resp *types.BlogTagOptionsResp, err error) {
	limit := int64(1000)
	if req != nil && req.Limit > 0 {
		limit = req.Limit
	}
	list, err := l.svcCtx.Domain.Blog.Tag.FindEnabledList(l.ctx, limit)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询标签选项失败", err)
	}
	items := make([]types.BlogTagOptionItem, 0, len(list))
	for _, t := range list {
		items = append(items, types.BlogTagOptionItem{
			Id:   t.Id,
			Name: t.Name,
		})
	}
	return &types.BlogTagOptionsResp{List: items}, nil
}
