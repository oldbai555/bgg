package logic

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/services/task/internal/consts"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
	"postapocgame/admin-server/services/task/internal/svc"
	"postapocgame/admin-server/services/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubmitTask 对应 AsyncTaskBackend.Submit 的 RPC 化，供 iam-rpc/sdk-rpc（当前阶段：单体
// 主进程）创建异步任务记录。原来的行为是"创建 admin_task 记录失败直接把错误返回给前端"，
// 是同步 RPC，不是 Streams（见 17-async-eventing.md 判断规则）。
func (l *SubmitTaskLogic) SubmitTask(in *task.SubmitTaskRequest) (*task.SubmitTaskResponse, error) {
	now := time.Now().Unix()
	t := &taskmodel.AdminTask{
		Name:          in.Name,
		Type:          in.TaskType,
		ExecutionType: in.ExecutionType,
		Status:        consts.TaskStatusPending,
		Params:        sql.NullString{String: in.Params, Valid: in.Params != ""},
		UserId:        in.UserId,
		ScheduledAt:   in.ScheduledAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	taskId, err := l.svcCtx.TaskRepo.Create(l.ctx, t)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	l.Infof("任务已创建: taskId=%d, userId=%d, name=%s", taskId, in.UserId, in.Name)

	return &task.SubmitTaskResponse{TaskId: taskId}, nil
}
