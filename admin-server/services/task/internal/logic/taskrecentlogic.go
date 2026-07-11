package logic

import (
	"context"

	"postapocgame/admin-server/services/task/internal/svc"
	"postapocgame/admin-server/services/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskRecentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskRecentLogic {
	return &TaskRecentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TaskRecent limit=0 时用 Config.RecentTaskLimit 兜底（静态配置，见 15-service-boundaries.md
// 第 5 节末尾建议，取代原来查 BusinessCache/字典的链路）。
func (l *TaskRecentLogic) TaskRecent(in *task.TaskRecentRequest) (*task.TaskRecentResponse, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = l.svcCtx.Config.RecentTaskLimit
	}
	if limit <= 0 {
		limit = 10
	}

	tasks, err := l.svcCtx.TaskRepo.FindRecent(l.ctx, limit, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	list := make([]*task.TaskItem, 0, len(tasks))
	for _, t := range tasks {
		item := &task.TaskItem{
			Id:            t.Id,
			Name:          t.Name,
			TaskType:      t.Type,
			ExecutionType: t.ExecutionType,
			Status:        t.Status,
			UserId:        t.UserId,
			ScheduledAt:   t.ScheduledAt,
			StartedAt:     t.StartedAt,
			FinishedAt:    t.FinishedAt,
			CreatedAt:     t.CreatedAt,
			ErrorMessage:  t.ErrorMessage,
		}
		if t.Params.Valid {
			item.Params = t.Params.String
		}
		if t.Result.Valid {
			item.Result = t.Result.String
		}
		list = append(list, item)
	}

	return &task.TaskRecentResponse{List: list}, nil
}
