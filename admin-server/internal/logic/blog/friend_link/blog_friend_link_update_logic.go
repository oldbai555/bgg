// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend_link

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkUpdateLogic {
	return &BlogFriendLinkUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogFriendLinkUpdate 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogfriendlinkupdatelogic.go。
func (l *BlogFriendLinkUpdateLogic) BlogFriendLinkUpdate(req *types.BlogFriendLinkUpdateReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogFriendLinkUpdate(l.ctx, &contentclient.BlogFriendLinkUpdateRequest{
		Id:       req.Id,
		Name:     req.Name,
		Url:      req.Url,
		Remark:   req.Remark,
		Status:   req.Status,
		OrderNum: req.OrderNum,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("更新友情链接失败", err)
	}
	return &types.Response{Code: 0, Message: "更新成功"}, nil
}
