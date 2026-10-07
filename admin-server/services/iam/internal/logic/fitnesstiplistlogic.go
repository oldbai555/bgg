package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipListLogic {
	return &FitnessTipListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTipListLogic) FitnessTipList(in *iam.FitnessTipListRequest) (*iam.FitnessTipListResponse, error) {
	list, total, err := l.svcCtx.Domain.Fitness.Tip.FindPage(l.ctx, in.Page, in.PageSize, in.Keyword)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询提示列表失败", err))
	}
	resp := &iam.FitnessTipListResponse{Total: total}
	for _, t := range list {
		resp.List = append(resp.List, fitnessTipItem(t))
	}
	return resp, nil
}
