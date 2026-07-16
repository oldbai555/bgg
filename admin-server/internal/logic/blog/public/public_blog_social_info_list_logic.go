// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

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

// PublicBlogSocialInfoList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicblogsocialinfolistlogic.go。
func (l *PublicBlogSocialInfoListLogic) PublicBlogSocialInfoList() (resp *types.PublicBlogSocialInfoListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogSocialInfoList(l.ctx, &contentclient.PublicBlogGlobalRequest{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询社交信息列表失败", err)
	}
	items := make([]types.PublicBlogSocialInfoItem, 0, len(rpcResp.List))
	for _, info := range rpcResp.List {
		items = append(items, types.PublicBlogSocialInfoItem{
			Id:       info.Id,
			Name:     info.Name,
			Url:      info.Url,
			Remark:   info.Remark,
			OrderNum: info.OrderNum,
		})
	}
	return &types.PublicBlogSocialInfoListResp{List: items}, nil
}
