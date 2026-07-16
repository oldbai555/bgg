package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/initdata"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDeleteLogic {
	return &RoleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleDeleteLogic) RoleDelete(in *iam.RoleDeleteRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "角色ID不能为空"))
	}
	if initdata.IsInitRoleID(in.Id) {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "初始化数据不可删除"))
	}

	if err := l.svcCtx.Domain.IAM.Role.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除角色失败", err))
	}
	return &iam.Empty{}, nil
}
