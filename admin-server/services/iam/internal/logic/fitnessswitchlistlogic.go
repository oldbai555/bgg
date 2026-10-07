package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessSwitchListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessSwitchListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessSwitchListLogic {
	return &FitnessSwitchListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessSwitchListLogic) FitnessSwitchList(in *iam.FitnessSwitchListRequest) (*iam.FitnessSwitchListResponse, error) {
	list, total, err := l.svcCtx.Domain.Fitness.SwitchViews(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessSwitchListResponse{Total: total}
	for _, v := range list {
		resp.List = append(resp.List, fitnessSwitchItem(v))
	}
	return resp, nil
}
