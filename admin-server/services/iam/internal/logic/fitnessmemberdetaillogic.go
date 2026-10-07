package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMemberDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMemberDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberDetailLogic {
	return &FitnessMemberDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMemberDetailLogic) FitnessMemberDetail(in *iam.FitnessMemberDetailRequest) (*iam.FitnessMemberDetailResponse, error) {
	detail, err := l.svcCtx.Domain.Fitness.MemberDetail(l.ctx, in.UserId, in.StartDate, in.EndDate)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	tplNames, err := l.svcCtx.Domain.Fitness.TemplateNames(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessMemberDetailResponse{
		Member:      fitnessMemberItem(detail.Member, detail.Summary, tplNames),
		RangeRate:   detail.RangeRate,
		BodyRecords: fitnessBodyItems(detail.Body),
	}
	for _, d := range detail.Calendar {
		resp.Calendar = append(resp.Calendar, &iam.FitnessCalendarDay{
			Date: d.Date, Counted: d.Counted, Checked: d.Checked, Score: d.Score, TrainingType: d.TrainingType,
		})
	}
	for _, v := range detail.Switches {
		resp.Switches = append(resp.Switches, fitnessSwitchItem(v))
	}
	return resp, nil
}
