package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	fitnessrepo "postapocgame/admin-server/services/iam/internal/repository/fitness"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessCheckinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessCheckinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessCheckinListLogic {
	return &FitnessCheckinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessCheckinListLogic) FitnessCheckinList(in *iam.FitnessCheckinListRequest) (*iam.FitnessCheckinListResponse, error) {
	filter := fitnessrepo.CheckinFilter{UserID: in.UserId, StartDate: in.StartDate, EndDate: in.EndDate}
	rows, total, err := l.svcCtx.Domain.Fitness.CheckinRows(l.ctx, filter, in.Page, in.PageSize)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessCheckinListResponse{Total: total}
	for _, r := range rows {
		c := r.Checkin
		resp.List = append(resp.List, &iam.FitnessCheckinRow{
			Id: c.Id, Date: c.CheckinDate, UserId: c.UserId, Nickname: r.Nickname, TemplateName: r.TemplateName,
			TrainingType: c.TrainingType, TrainingStatus: c.TrainingStatus,
			ExerciseDone: r.ExerciseDone, ExerciseTotal: c.ExerciseTotal,
			MealOnPlan: r.MealOnPlan, MealOther: r.MealOther, MealSkipped: r.MealSkipped,
			Steps: c.Steps, StepGoal: c.StepGoal, WaterCups: c.WaterCups, Score: r.Score, UpdatedAt: c.UpdatedAt,
		})
	}
	return resp, nil
}
