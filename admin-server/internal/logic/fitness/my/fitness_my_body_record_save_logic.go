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

type FitnessMyBodyRecordSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyBodyRecordSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyBodyRecordSaveLogic {
	return &FitnessMyBodyRecordSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyBodyRecordSaveLogic) FitnessMyBodyRecordSave(req *types.FitnessMyBodyRecordSaveReq) error {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err := l.svcCtx.IamRPC.FitnessMyBodyRecordSave(l.ctx, &iamclient.FitnessMyBodyRecordSaveRequest{
		UserId:   user.UserID,
		Date:     req.Date,
		WeightKg: req.WeightKg,
		WaistCm:  req.WaistCm,
	})
	if err != nil {
		return errs.WrapGRPCError("保存身体数据失败", err)
	}
	return nil
}
