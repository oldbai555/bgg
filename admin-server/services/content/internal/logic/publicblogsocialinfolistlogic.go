package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogSocialInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogSocialInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogSocialInfoListLogic {
	return &PublicBlogSocialInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogSocialInfoList 迁移自
// internal/logic/blog/public/public_blog_social_info_list_logic.go。
func (l *PublicBlogSocialInfoListLogic) PublicBlogSocialInfoList(in *content.PublicBlogGlobalRequest) (*content.PublicBlogSocialInfoListResponse, error) {
	list, err := l.svcCtx.SocialInfo.FindEnabledList(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询社交信息列表失败", err))
	}

	items := make([]*content.PublicBlogSocialInfoItem, 0, len(list))
	for _, info := range list {
		items = append(items, &content.PublicBlogSocialInfoItem{
			Id:       info.Id,
			Name:     info.Name,
			Url:      info.Url,
			Remark:   info.Remark,
			OrderNum: info.OrderNum,
		})
	}

	return &content.PublicBlogSocialInfoListResponse{List: items}, nil
}
