// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	taskrepo "postapocgame/admin-server/internal/repository/task"
)

type TaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskListLogic) TaskList(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	// 参数默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 构建查询过滤条件
	filters := &taskrepo.TaskQueryFilter{
		Name:          req.Name,
		Type:          req.TaskType,
		ExecutionType: req.ExecutionType,
		Status:        req.Status,
		UserId:        req.UserId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	}

	// 查询任务列表
	taskRepo := taskrepo.NewTaskRepository(l.svcCtx.Repository)
	tasks, total, err := taskRepo.FindPage(l.ctx, page, pageSize, filters)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询任务列表失败", err)
	}

	// 转换为响应格式
	list := make([]types.TaskItem, 0, len(tasks))
	for _, task := range tasks {
		params := ""
		if task.Params.Valid {
			params = task.Params.String
		}
		result := ""
		if task.Result.Valid {
			result = task.Result.String
		}

		list = append(list, types.TaskItem{
			Id:            task.Id,
			Name:          task.Name,
			TaskType:      task.Type,
			ExecutionType: task.ExecutionType,
			Status:        task.Status,
			UserId:        task.UserId,
			ScheduledAt:   task.ScheduledAt,
			StartedAt:     task.StartedAt,
			FinishedAt:    task.FinishedAt,
			CreatedAt:     task.CreatedAt,
			Params:        params,
			Result:        result,
			ErrorMessage:  task.ErrorMessage,
		})
	}

	return &types.TaskListResp{
		Total: total,
		List:  list,
	}, nil
}
