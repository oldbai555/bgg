// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_friend_link

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *BlogFriendLinkDeleteLogic) BlogFriendLinkDelete(req *types.BlogFriendLinkDeleteReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "链接ID不能为空")
	}

	// 检查记录是否存在
	link, err := l.svcCtx.BlogFriendLinkRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询友情链接失败", err)
	}
	if link == nil {
		return nil, errs.New(errs.CodeNotFound, "友情链接不存在")
	}

	// 执行软删除
	if err := l.svcCtx.BlogFriendLinkRepository.Delete(l.ctx, req.Id); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "删除成功",
	}, nil
}
