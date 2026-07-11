// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

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

func (l *TaskCancelLogic) TaskCancel(req *types.TaskCancelReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "任务ID不能为空")
	}

	// 获取当前用户
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	// 查询任务
	task, err := l.svcCtx.Domain.Task.Task.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 验证权限：只能取消自己创建的任务（或管理员）
	if task.UserId != user.UserID && user.UserID != 1 {
		return nil, errs.New(errs.CodeForbidden, "无权取消该任务")
	}

	// 只能取消未开始或进行中的任务
	if task.Status != consts.TaskStatusPending && task.Status != consts.TaskStatusRunning {
		return nil, errs.New(errs.CodeBadRequest, "只能取消未开始或进行中的任务")
	}

	// 更新任务状态为失败（取消视为失败）
	// 注意：这里不更新finished_at，因为任务是被取消的，不是正常完成的
	err = l.svcCtx.Domain.Task.Task.UpdateStatus(l.ctx, req.Id, consts.TaskStatusFailed, 0, 0)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "取消任务失败", err)
	}

	logx.Infof("任务已取消: taskId=%d, userId=%d", req.Id, user.UserID)

	return &types.Response{
		Code:    0,
		Message: "任务已取消",
	}, nil
}
