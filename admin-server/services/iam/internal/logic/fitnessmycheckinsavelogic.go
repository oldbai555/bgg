package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	fitnessdomain "postapocgame/admin-server/services/iam/internal/domain/fitness"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyCheckinSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyCheckinSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyCheckinSaveLogic {
	return &FitnessMyCheckinSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyCheckinSaveLogic) FitnessMyCheckinSave(in *iam.FitnessMyCheckinSaveRequest) (*iam.FitnessCheckinItem, error) {
	input := fitnessdomain.CheckinInput{
		Date:           in.Date,
		TrainingStatus: in.TrainingStatus,
		ExerciseDone:   in.ExerciseDone,
		Steps:          in.Steps,
		WaterCups:      in.WaterCups,
	}
	for _, m := range in.Meals {
		input.Meals = append(input.Meals, fitnessdomain.MealCheck{Slot: m.Slot, Status: m.Status, Note: m.Note})
	}
	saved, err := l.svcCtx.Domain.Fitness.SaveCheckin(l.ctx, in.UserId, input)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return fitnessCheckinItem(saved), nil
}
