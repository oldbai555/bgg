// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"context"

	"postapocgame/admin-server/internal/logic/fitness/fitnessconv"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyCheckinSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyCheckinSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyCheckinSaveLogic {
	return &FitnessMyCheckinSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyCheckinSaveLogic) FitnessMyCheckinSave(req *types.FitnessMyCheckinSaveReq) (resp *types.FitnessCheckinItem, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	meals := make([]*iamclient.FitnessMealCheck, 0, len(req.Meals))
	for _, m := range req.Meals {
		meals = append(meals, &iamclient.FitnessMealCheck{Slot: m.Slot, Status: m.Status, Note: m.Note})
	}
	rpcResp, err := l.svcCtx.IamRPC.FitnessMyCheckinSave(l.ctx, &iamclient.FitnessMyCheckinSaveRequest{
		UserId:         user.UserID,
		Date:           req.Date,
		TrainingStatus: req.TrainingStatus,
		ExerciseDone:   req.ExerciseDone,
		Meals:          meals,
		Steps:          req.Steps,
		WaterCups:      req.WaterCups,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("保存打卡失败", err)
	}
	return fitnessconv.CheckinItem(rpcResp), nil
}
