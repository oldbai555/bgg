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

type BlogFriendLinkDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkDeleteLogic {
	return &BlogFriendLinkDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogFriendLinkDelete 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogfriendlinkdeletelogic.go。
func (l *BlogFriendLinkDeleteLogic) BlogFriendLinkDelete(req *types.BlogFriendLinkDeleteReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogFriendLinkDelete(l.ctx, &contentclient.BlogFriendLinkDeleteRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("删除友情链接失败", err)
	}
	return &types.Response{Code: 0, Message: "删除成功"}, nil
}
