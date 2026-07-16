package logic

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"postapocgame/admin-server/services/task/internal/consts"
	"postapocgame/admin-server/services/task/internal/svc"
	"postapocgame/admin-server/services/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCancelLogic {
	return &TaskCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TaskCancel 权限归属/状态校验从 gateway 侧下沉到这里（task-rpc 拆分后不再能访问登录态，
// gateway 显式传 operator_user_id，见 services/task/rpc/task.proto 的字段注释）。
func (l *TaskCancelLogic) TaskCancel(in *task.TaskCancelRequest) (*task.TaskCancelResponse, error) {
	t, err := l.svcCtx.TaskRepo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	// 验证权限：只能取消自己创建的任务（或管理员）
	if t.UserId != in.OperatorUserId && in.OperatorUserId != consts.SuperAdminUserID {
		return nil, status.Error(codes.PermissionDenied, "无权取消该任务")
	}

	// 只能取消未开始或进行中的任务
	if t.Status != consts.TaskStatusPending && t.Status != consts.TaskStatusRunning {
		return nil, status.Error(codes.FailedPrecondition, "只能取消未开始或进行中的任务")
	}

	// 更新任务状态为失败（取消视为失败）；不更新 finished_at，因为任务是被取消的，不是正常完成的
	if err := l.svcCtx.TaskRepo.UpdateStatus(l.ctx, in.Id, consts.TaskStatusFailed, 0, 0); err != nil {
		return nil, toGRPCStatus(err)
	}

	l.Infof("任务已取消: taskId=%d, operatorUserId=%d", in.Id, in.OperatorUserId)

	return &task.TaskCancelResponse{}, nil
}
