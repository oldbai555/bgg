package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictTypeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeCreateLogic {
	return &DictTypeCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictTypeCreateLogic) DictTypeCreate(in *iam.DictTypeCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" || in.Code == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型名称和编码不能为空"))
	}

	if _, err := l.svcCtx.Domain.System.DictType.FindByCode(l.ctx, in.Code); err == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型编码已存在"))
	}

	status := in.Status
	if status == 0 {
		status = 1
	}

	dictType := systemmodel.AdminDictType{
		Name:        in.Name,
		Code:        in.Code,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		Status:      status,
	}

	if err := l.svcCtx.Domain.System.DictType.Create(l.ctx, &dictType); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建字典类型失败", err))
	}
	return &iam.Empty{}, nil
}
