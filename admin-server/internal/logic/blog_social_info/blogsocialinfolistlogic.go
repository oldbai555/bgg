// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_social_info

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogSocialInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoListLogic {
	return &BlogSocialInfoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogSocialInfoListLogic) BlogSocialInfoList(req *types.BlogSocialInfoListReq) (resp *types.BlogSocialInfoListResp, err error) {
	// 参数预处理与默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	// 调用仓储层分页查询
	list, total, err := l.svcCtx.BlogSocialInfoRepository.FindPage(l.ctx, page, pageSize, req.Status, req.Keyword)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询社交信息列表失败", err)
	}

	items := make([]types.BlogSocialInfoItem, 0, len(list))
	for _, info := range list {
		items = append(items, types.BlogSocialInfoItem{
			Id:        info.Id,
			Name:      info.Name,
			Url:       info.Url,
			Remark:    info.Remark,
			Status:    info.Status,
			OrderNum:  info.OrderNum,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
		})
	}

	return &types.BlogSocialInfoListResp{
		Page:  page,
		Size:  pageSize,
		Total: total,
		List:  items,
	}, nil
}
