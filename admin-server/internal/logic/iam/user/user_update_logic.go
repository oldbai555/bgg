// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "用户ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.UserUpdate(l.ctx, &iamclient.UserUpdateRequest{
		Id:           req.Id,
		Username:     req.Username,
		Nickname:     req.Nickname,
		Password:     req.Password,
		Avatar:       req.Avatar,
		Signature:    req.Signature,
		DepartmentId: req.DepartmentId,
		Status:       req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新用户失败", err)
	}
	return nil
}
