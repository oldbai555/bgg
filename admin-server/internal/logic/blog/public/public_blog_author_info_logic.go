// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogAuthorInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogAuthorInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogAuthorInfoLogic {
	return &PublicBlogAuthorInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicBlogAuthorInfoLogic) PublicBlogAuthorInfo() (resp *types.PublicBlogAuthorInfoResp, err error) {
	// 查询超级管理员（id=1）的信息
	user, err := iamrepo.NewUserRepository(l.svcCtx.Repository).FindByID(l.ctx, 1)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询作者信息失败", err)
	}

	// 如果用户不存在或已删除，返回默认值
	if user == nil || user.DeletedAt > 0 {
		return &types.PublicBlogAuthorInfoResp{
			Id:        1,
			Nickname:  "管理员",
			Avatar:    "",
			Signature: "",
		}, nil
	}

	return &types.PublicBlogAuthorInfoResp{
		Id:        user.Id,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Signature: user.Signature,
	}, nil
}
