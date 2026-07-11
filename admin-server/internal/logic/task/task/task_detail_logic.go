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

	rpcResp, err := l.svcCtx.TaskRPC.TaskDetail(l.ctx, &taskclient.TaskDetailRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询任务详情失败", err)
	}

	return &types.TaskDetailResp{
		Id:            rpcResp.Id,
		Name:          rpcResp.Name,
		TaskType:      rpcResp.TaskType,
		ExecutionType: rpcResp.ExecutionType,
		Status:        rpcResp.Status,
		Params:        rpcResp.Params,
		Result:        rpcResp.Result,
		ErrorMessage:  rpcResp.ErrorMessage,
		UserId:        rpcResp.UserId,
		ScheduledAt:   rpcResp.ScheduledAt,
		StartedAt:     rpcResp.StartedAt,
		FinishedAt:    rpcResp.FinishedAt,
		CreatedAt:     rpcResp.CreatedAt,
		UpdatedAt:     rpcResp.UpdatedAt,
	}, nil
}
