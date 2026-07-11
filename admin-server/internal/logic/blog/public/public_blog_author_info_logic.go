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
	// TODO(phase2-content-rpc): 跨域读取 IAM 用户信息，Phase 2 拆分后改为调用 iam-rpc.GetUserProfile
	user, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, 1)
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
