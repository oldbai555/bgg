package logic

import (
	"context"

	"postapocgame/admin-server/services/task/internal/repository"
	"postapocgame/admin-server/services/task/internal/svc"
	"postapocgame/admin-server/services/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskListLogic) TaskList(in *task.TaskListRequest) (*task.TaskListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	filters := &repository.TaskQueryFilter{
		Name:          in.Name,
		Type:          in.TaskType,
		ExecutionType: in.ExecutionType,
		Status:        in.Status,
		UserId:        in.UserId,
		StartTime:     in.StartTime,
		EndTime:       in.EndTime,
	}

	tasks, total, err := l.svcCtx.TaskRepo.FindPage(l.ctx, page, pageSize, filters)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	list := make([]*task.TaskItem, 0, len(tasks))
	for _, t := range tasks {
		params := ""
		if t.Params.Valid {
			params = t.Params.String
		}
		result := ""
		if t.Result.Valid {
			result = t.Result.String
		}
		list = append(list, &task.TaskItem{
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
			Params:        params,
			Result:        result,
			ErrorMessage:  t.ErrorMessage,
		})
	}

	return &task.TaskListResponse{
		Total: total,
		List:  list,
	}, nil
}
