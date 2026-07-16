package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictTypeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeUpdateLogic {
	return &DictTypeUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictTypeUpdateLogic) DictTypeUpdate(in *iam.DictTypeUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型ID不能为空"))
	}

	dictType, err := l.svcCtx.Domain.System.DictType.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典类型失败", err))
	}

	if in.Name != "" {
		dictType.Name = in.Name
	}
	if in.Description != "" {
		dictType.Description = sql.NullString{String: in.Description, Valid: true}
	}
	if in.Status == 0 || in.Status == 1 {
		dictType.Status = in.Status
	}

	if err := l.svcCtx.Domain.System.DictType.Update(l.ctx, dictType); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新字典类型失败", err))
	}
	return &iam.Empty{}, nil
}
