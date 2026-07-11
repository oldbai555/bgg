// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/core/logx"
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

// TaskList 薄胶水：解析 HTTP 请求 -> 拼一次 TaskRPC 请求 -> 映射响应，task 域的实际
// 业务逻辑已经搬进 services/task/internal/logic/tasklistlogic.go。
func (l *TaskListLogic) TaskList(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	rpcResp, err := l.svcCtx.TaskRPC.TaskList(l.ctx, &taskclient.TaskListRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		Name:          req.Name,
		TaskType:      req.TaskType,
		ExecutionType: req.ExecutionType,
		Status:        req.Status,
		UserId:        req.UserId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询任务列表失败", err)
	}

	list := make([]types.TaskItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.TaskItem{
			Id:            item.Id,
			Name:          item.Name,
			TaskType:      item.TaskType,
			ExecutionType: item.ExecutionType,
			Status:        item.Status,
			UserId:        item.UserId,
			ScheduledAt:   item.ScheduledAt,
			StartedAt:     item.StartedAt,
			FinishedAt:    item.FinishedAt,
			CreatedAt:     item.CreatedAt,
			Params:        item.Params,
			Result:        item.Result,
			ErrorMessage:  item.ErrorMessage,
		})
	}

	return &types.TaskListResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
