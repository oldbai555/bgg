package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyBodyRecordListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyBodyRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyBodyRecordListLogic {
	return &FitnessMyBodyRecordListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyBodyRecordListLogic) FitnessMyBodyRecordList(in *iam.FitnessMyBodyRecordListRequest) (*iam.FitnessBodyRecordListResponse, error) {
	points, err := l.svcCtx.Domain.Fitness.MyBodyPoints(l.ctx, in.UserId, in.Days)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.FitnessBodyRecordListResponse{List: fitnessBodyItems(points)}, nil
}
