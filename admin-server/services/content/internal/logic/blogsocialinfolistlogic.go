package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogSocialInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoListLogic {
	return &BlogSocialInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 社交信息

// BlogSocialInfoList 迁移自 internal/logic/blog/social_info/blog_social_info_list_logic.go。
func (l *BlogSocialInfoListLogic) BlogSocialInfoList(in *content.BlogSocialInfoListRequest) (*content.BlogSocialInfoListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.SocialInfo.FindPage(l.ctx, page, pageSize, in.Status, in.Keyword)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询社交信息列表失败", err))
	}

	items := make([]*content.BlogSocialInfoItem, 0, len(list))
	for _, info := range list {
		items = append(items, &content.BlogSocialInfoItem{
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

	return &content.BlogSocialInfoListResponse{Page: page, Size: pageSize, Total: total, List: items}, nil
}
