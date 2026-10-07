package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMySuggestionDecideLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMySuggestionDecideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMySuggestionDecideLogic {
	return &FitnessMySuggestionDecideLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMySuggestionDecideLogic) FitnessMySuggestionDecide(in *iam.FitnessMySuggestionDecideRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.DecideSuggestion(l.ctx, in.UserId, in.Id, in.Accept); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
