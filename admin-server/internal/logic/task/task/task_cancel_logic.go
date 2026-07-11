// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCancelLogic {
	return &TaskCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TaskCancel 薄胶水：task-rpc 拆分后不再能访问登录态，权限归属/状态校验都下沉到
// services/task/internal/logic/taskcancellogic.go，这里只负责从 JWT context 取出
// operator_user_id 显式传过去。
func (l *TaskCancelLogic) TaskCancel(req *types.TaskCancelReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "任务ID不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	if _, err := l.svcCtx.TaskRPC.TaskCancel(l.ctx, &taskclient.TaskCancelRequest{
		Id:             req.Id,
		OperatorUserId: user.UserID,
	}); err != nil {
		return nil, errs.WrapGRPCError("取消任务失败", err)
	}

	logx.Infof("任务已取消: taskId=%d, userId=%d", req.Id, user.UserID)

	return &types.Response{
		Code:    0,
		Message: "任务已取消",
	}, nil
}
