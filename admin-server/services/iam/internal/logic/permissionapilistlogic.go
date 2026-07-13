package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionApiListLogic {
	return &PermissionApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionApiListLogic) PermissionApiList(in *iam.PermissionApiListRequest) (*iam.PermissionApiListResponse, error) {
	if in.PermissionId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}

	if _, err := l.svcCtx.Domain.IAM.Permission.FindByID(l.ctx, in.PermissionId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "权限不存在", err))
	}

	apiIDs, err := l.svcCtx.Domain.IAM.PermissionApi.ListApiIDsByPermissionID(l.ctx, in.PermissionId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询权限接口失败", err))
	}

	return &iam.PermissionApiListResponse{ApiIds: apiIDs}, nil
}
