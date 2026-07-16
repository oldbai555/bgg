package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogFriendLinkDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkDeleteLogic {
	return &BlogFriendLinkDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogFriendLinkDelete 迁移自 internal/logic/blog/friend_link/blog_friend_link_delete_logic.go。
func (l *BlogFriendLinkDeleteLogic) BlogFriendLinkDelete(in *content.BlogFriendLinkDeleteRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "链接ID不能为空"))
	}
	link, err := l.svcCtx.FriendLink.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询友情链接失败", err))
	}
	if link == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "友情链接不存在"))
	}
	if err := l.svcCtx.FriendLink.Delete(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &content.Empty{}, nil
}
