// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessCheckinExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessCheckinExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessCheckinExportLogic {
	return &FitnessCheckinExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessCheckinExportLogic) FitnessCheckinExport(req *types.FitnessCheckinExportReq) (resp *types.Response, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	filters := make(map[string]interface{})
	if req.UserId > 0 {
		filters[consts.TaskFilterUserId] = req.UserId
	}
	if req.StartDate != "" {
		filters[consts.TaskFilterStartDate] = req.StartDate
	}
	if req.EndDate != "" {
		filters[consts.TaskFilterEndDate] = req.EndDate
	}
	paramsJSON, err := json.Marshal(map[string]interface{}{
		"module":  consts.TaskModuleFitnessCheckin,
		"filters": filters,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "序列化任务参数失败", err)
	}

	rpcResp, err := l.svcCtx.TaskRPC.SubmitTask(l.ctx, &taskclient.SubmitTaskRequest{
		Name:          fmt.Sprintf("身材管理打卡明细导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		TaskType:      consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Params:        string(paramsJSON),
		UserId:        user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建导出任务失败", err)
	}
	logx.Infof("身材管理打卡明细导出任务已创建: taskId=%d, userId=%d", rpcResp.TaskId, user.UserID)

	return &types.Response{
		Code:    errs.CodeOK,
		Message: "已创建异步导出任务，请在右下角任务列表查看进度",
	}, nil
}
