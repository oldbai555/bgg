package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberListLogic {
	return &FitnessMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMemberListLogic) FitnessMemberList(in *iam.FitnessMemberListRequest) (*iam.FitnessMemberListResponse, error) {
	members, summaries, total, err := l.svcCtx.Domain.Fitness.MemberPage(l.ctx, in.Page, in.PageSize, in.Keyword)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	tplNames, err := l.svcCtx.Domain.Fitness.TemplateNames(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessMemberListResponse{Total: total}
	for _, m := range members {
		resp.List = append(resp.List, fitnessMemberItem(m, summaries[m.UserId], tplNames))
	}
	return resp, nil
}
