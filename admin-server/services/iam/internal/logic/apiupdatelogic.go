package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdateLogic {
	return &ApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiUpdateLogic) ApiUpdate(in *iam.ApiUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "接口ID不能为空"))
	}

	api, err := l.svcCtx.Domain.IAM.Api.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询接口失败", err))
	}

	if in.Name != "" {
		api.Name = in.Name
	}
	if in.Method != "" {
		api.Method = in.Method
	}
	if in.Path != "" {
		api.Path = in.Path
	}
	if in.Description != "" {
		api.Description = sql.NullString{String: in.Description, Valid: true}
	}
	if in.Status == 0 || in.Status == 1 {
		api.Status = in.Status
	}

	if err := l.svcCtx.Domain.IAM.Api.Update(l.ctx, api); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新接口失败", err))
	}
	return &iam.Empty{}, nil
}
