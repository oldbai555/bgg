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

type TaskDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDetailLogic {
	return &TaskDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskDetailLogic) TaskDetail(req *types.TaskDetailReq) (resp *types.TaskDetailResp, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "任务ID不能为空")
	}

	// 查询任务详情
	taskRepo := taskrepo.NewTaskRepository(l.svcCtx.Repository)
	task, err := taskRepo.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	resp = &types.TaskDetailResp{
		Id:            task.Id,
		Name:          task.Name,
		TaskType:      task.Type,
		ExecutionType: task.ExecutionType,
		Status:        task.Status,
		Params:        "",
		Result:        "",
		ErrorMessage:  task.ErrorMessage,
		UserId:        task.UserId,
		ScheduledAt:   task.ScheduledAt,
		StartedAt:     task.StartedAt,
		FinishedAt:    task.FinishedAt,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}

	// 处理Params和Result（可能为NULL）
	if task.Params.Valid {
		resp.Params = task.Params.String
	}
	if task.Result.Valid {
		resp.Result = task.Result.String
	}

	return resp, nil
}
