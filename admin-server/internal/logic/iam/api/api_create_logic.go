// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"postapocgame/admin-server/internal/model/iam"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCreateLogic {
	return &ApiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCreateLogic) ApiCreate(req *types.ApiCreateReq) error {
	if req == nil || req.Name == "" || req.Method == "" || req.Path == "" {
		return errs.New(errs.CodeBadRequest, "接口名称、方法和路径不能为空")
	}

	// 检查是否已存在相同的 method+path
	_, err := l.svcCtx.Domain.IAM.Api.FindByMethodAndPath(l.ctx, req.Method, req.Path)
	if err == nil {
		return errs.New(errs.CodeBadRequest, "该接口已存在")
	}

	api := iam.AdminApi{
		Name:        req.Name,
		Method:      req.Method,
		Path:        req.Path,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Status:      req.Status,
	}
	if api.Status == 0 {
		api.Status = 1
	}

	if err := l.svcCtx.Domain.IAM.Api.Create(l.ctx, &api); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建接口失败", err)
	}
	return nil
}
