package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateListLogic {
	return &FitnessTemplateListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Fitness 后台：计划模板 / 通用提示 / 使用者 / 打卡统计（导出走 task-rpc → TaskCallback.FetchExportData）
func (l *FitnessTemplateListLogic) FitnessTemplateList(in *iam.FitnessTemplateListRequest) (*iam.FitnessTemplateListResponse, error) {
	list, total, err := l.svcCtx.Domain.Fitness.Template.FindPage(l.ctx, in.Page, in.PageSize, in.Keyword)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询模板列表失败", err))
	}
	resp := &iam.FitnessTemplateListResponse{Total: total}
	for _, t := range list {
		resp.List = append(resp.List, fitnessTemplateItem(t))
	}
	return resp, nil
}
