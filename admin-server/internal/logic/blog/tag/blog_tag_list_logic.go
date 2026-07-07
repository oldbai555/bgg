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

func (l *BlogTagListLogic) BlogTagList(req *types.BlogTagListReq) (resp *types.BlogTagListResp, err error) {
	// 参数预处理与默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 调用仓储层分页查询
	tagRepo := blogrepo.NewBlogTagRepository(l.svcCtx.Repository)
	list, total, err := tagRepo.FindPage(l.ctx, page, pageSize, req.Name, req.Status)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询标签列表失败", err)
	}

	items := make([]types.BlogTagItem, 0, len(list))
	for _, t := range list {
		items = append(items, types.BlogTagItem{
			Id:        t.Id,
			Name:      t.Name,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return &types.BlogTagListResp{
		Total: total,
		List:  items,
	}, nil
}
