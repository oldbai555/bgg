package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserListLogic) UserList(in *iam.UserListRequest) (*iam.UserListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.IAM.User.FindPage(l.ctx, in.Page, in.PageSize, in.Username)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户列表失败", err))
	}

	items := make([]*iam.UserItem, 0, len(list))
	for _, u := range list {
		items = append(items, &iam.UserItem{
			Id:           u.Id,
			Username:     u.Username,
			Nickname:     u.Nickname,
			Avatar:       u.Avatar,
			Signature:    u.Signature,
			DepartmentId: u.DepartmentId,
			Status:       u.Status,
			CreatedAt:    int64(u.CreatedAt),
		})
	}

	return &iam.UserListResponse{
		Total: total,
		List:  items,
	}, nil
}
