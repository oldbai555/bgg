// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogSocialInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogSocialInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogSocialInfoListLogic {
	return &PublicBlogSocialInfoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogSocialInfoListLogic) PublicBlogSocialInfoList() (resp *types.PublicBlogSocialInfoListResp, err error) {
	// 查询启用的社交信息列表
	list, err := l.svcCtx.Domain.Blog.SocialInfo.FindEnabledList(l.ctx)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询社交信息列表失败", err)
	}

	items := make([]types.PublicBlogSocialInfoItem, 0, len(list))
	for _, info := range list {
		items = append(items, types.PublicBlogSocialInfoItem{
			Id:       info.Id,
			Name:     info.Name,
			Url:      info.Url,
			Remark:   info.Remark,
			OrderNum: info.OrderNum,
		})
	}

	return &types.PublicBlogSocialInfoListResp{
		List: items,
	}, nil
}
