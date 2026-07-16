package logic

import (
	"context"

	"postapocgame/admin-server/services/task/internal/svc"
	"postapocgame/admin-server/services/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDetailLogic {
	return &TaskDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskDetailLogic) TaskDetail(in *task.TaskDetailRequest) (*task.TaskDetailResponse, error) {
	t, err := l.svcCtx.TaskRepo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	resp := &task.TaskDetailResponse{
		Id:            t.Id,
		Name:          t.Name,
		TaskType:      t.Type,
		ExecutionType: t.ExecutionType,
		Status:        t.Status,
		ErrorMessage:  t.ErrorMessage,
		UserId:        t.UserId,
		ScheduledAt:   t.ScheduledAt,
		StartedAt:     t.StartedAt,
		FinishedAt:    t.FinishedAt,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
	if t.Params.Valid {
		resp.Params = t.Params.String
	}
	if t.Result.Valid {
		resp.Result = t.Result.String
	}
	return resp, nil
}
