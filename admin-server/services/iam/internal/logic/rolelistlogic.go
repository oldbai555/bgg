package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleListLogic) RoleList(in *iam.RoleListRequest) (*iam.RoleListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.IAM.Role.FindPage(l.ctx, in.Page, in.PageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询角色列表失败", err))
	}

	items := make([]*iam.RoleItem, 0, len(list))
	for _, r := range list {
		description := ""
		if r.Description.Valid {
			description = r.Description.String
		}
		items = append(items, &iam.RoleItem{
			Id:          r.Id,
			Name:        r.Name,
			Code:        r.Code,
			Description: description,
			Status:      r.Status,
		})
	}

	return &iam.RoleListResponse{
		Total: total,
		List:  items,
	}, nil
}
