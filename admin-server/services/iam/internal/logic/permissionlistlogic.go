package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionListLogic {
	return &PermissionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionListLogic) PermissionList(in *iam.PermissionListRequest) (*iam.PermissionListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.IAM.Permission.FindPage(l.ctx, in.Page, in.PageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询权限列表失败", err))
	}

	items := make([]*iam.PermissionItem, 0, len(list))
	for _, p := range list {
		description := ""
		if p.Description.Valid {
			description = p.Description.String
		}
		items = append(items, &iam.PermissionItem{
			Id:          p.Id,
			Name:        p.Name,
			Code:        p.Code,
			Description: description,
		})
	}

	return &iam.PermissionListResponse{
		Total: total,
		List:  items,
	}, nil
}
