package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCreateLogic {
	return &ApiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiCreateLogic) ApiCreate(in *iam.ApiCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" || in.Method == "" || in.Path == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "接口名称、方法和路径不能为空"))
	}

	_, err := l.svcCtx.Domain.IAM.Api.FindByMethodAndPath(l.ctx, in.Method, in.Path)
	if err == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该接口已存在"))
	}
	if !isErrNotFound(err) {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询接口失败", err))
	}

	api := iammodel.AdminApi{
		Name:        in.Name,
		Method:      in.Method,
		Path:        in.Path,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		Status:      in.Status,
	}
	if api.Status == 0 {
		api.Status = 1
	}

	if err := l.svcCtx.Domain.IAM.Api.Create(l.ctx, &api); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建接口失败", err))
	}
	return &iam.Empty{}, nil
}
