// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package social_info

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

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

// BlogSocialInfoList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogsocialinfolistlogic.go。
func (l *BlogSocialInfoListLogic) BlogSocialInfoList(req *types.BlogSocialInfoListReq) (resp *types.BlogSocialInfoListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.BlogSocialInfoList(l.ctx, &contentclient.BlogSocialInfoListRequest{
		Page:    req.Page,
		Size:    req.Size,
		Status:  req.Status,
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询社交信息列表失败", err)
	}

	items := make([]types.BlogSocialInfoItem, 0, len(rpcResp.List))
	for _, info := range rpcResp.List {
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

	return &types.BlogSocialInfoListResp{Page: rpcResp.Page, Size: rpcResp.Size, Total: rpcResp.Total, List: items}, nil
}
