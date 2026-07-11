// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package performance_log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/task/taskclient"
)

type PerformanceLogExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPerformanceLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PerformanceLogExportLogic {
	return &PerformanceLogExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PerformanceLogExportLogic) PerformanceLogExport(req *types.PerformanceLogExportReq) (*types.Response, error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	filters := make(map[string]interface{})
	if req.Method != "" {
		filters[consts.TaskFilterMethod] = req.Method
	}
	if req.Path != "" {
		filters[consts.TaskFilterPath] = req.Path
	}
	if req.IsSlow != 0 {
		filters[consts.TaskFilterIsSlow] = req.IsSlow
	}
	if req.StatusCode != 0 {
		filters[consts.TaskFilterStatusCode] = req.StatusCode
	}
	if req.StartTime != "" {
		filters[consts.TaskFilterStartTime] = req.StartTime
	}
	if req.EndTime != "" {
		filters[consts.TaskFilterEndTime] = req.EndTime
	}

	paramsJSON, err := json.Marshal(map[string]interface{}{
		"module":  consts.TaskModulePerformanceLog,
		"filters": filters,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "序列化任务参数失败", err)
	}

	resp, err := l.svcCtx.TaskRPC.SubmitTask(l.ctx, &taskclient.SubmitTaskRequest{
		Name:          fmt.Sprintf("性能日志导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		TaskType:      consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Params:        string(paramsJSON),
		UserId:        user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建导出任务失败", err)
	}

	logx.Infof("性能日志导出任务已创建: taskId=%d, userId=%d", resp.TaskId, user.UserID)

	return &types.Response{
		Code:    0,
		Message: "已创建异步导出任务，请在右下角任务列表查看进度",
	}, nil
}
