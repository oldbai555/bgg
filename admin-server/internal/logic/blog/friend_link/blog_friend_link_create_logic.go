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

type BlogFriendLinkCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkCreateLogic {
	return &BlogFriendLinkCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogFriendLinkCreate 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogfriendlinkcreatelogic.go。
func (l *BlogFriendLinkCreateLogic) BlogFriendLinkCreate(req *types.BlogFriendLinkCreateReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogFriendLinkCreate(l.ctx, &contentclient.BlogFriendLinkCreateRequest{
		Name:     req.Name,
		Url:      req.Url,
		Remark:   req.Remark,
		Status:   req.Status,
		OrderNum: req.OrderNum,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建友情链接失败", err)
	}
	return &types.Response{Code: 0, Message: "创建成功"}, nil
}
