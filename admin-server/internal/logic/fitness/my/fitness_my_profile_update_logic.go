// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyProfileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyProfileUpdateLogic {
	return &FitnessMyProfileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyProfileUpdateLogic) FitnessMyProfileUpdate(req *types.FitnessMyProfileUpdateReq) error {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err := l.svcCtx.IamRPC.FitnessMyProfileUpdate(l.ctx, &iamclient.FitnessMyProfileUpdateRequest{
		UserId:             user.UserID,
		HeightCm:           req.HeightCm,
		TargetWeightKg:     req.TargetWeightKg,
		TrainingPreference: req.TrainingPreference,
	})
	if err != nil {
		return errs.WrapGRPCError("更新个人信息失败", err)
	}
	return nil
}
