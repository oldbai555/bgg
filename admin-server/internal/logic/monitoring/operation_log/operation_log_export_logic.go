// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package operation_log

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"postapocgame/admin-server/internal/domain/task"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

	taskmodel "postapocgame/admin-server/internal/model/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperationLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogExportLogic {
	return &OperationLogExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// OperationLogExport 将操作日志导出改为异步任务：
// 1. 根据筛选条件构造 ExcelExportParams
// 2. 创建异步任务记录（task_type=异步导出Excel，execution_type=异步）
// 3. 由任务调度器 + ExcelExportExecutor 实际生成文件，并在任务结果中写入下载URL
func (l *OperationLogExportLogic) OperationLogExport(req *types.OperationLogExportReq) (*types.OperationLogExportResp, error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 获取当前登录用户
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	// 构造导出参数（ExcelExportParams）
	filters := make(map[string]interface{})
	if req.UserId > 0 {
		filters["userId"] = req.UserId
	}
	if req.Username != "" {
		filters["username"] = req.Username
	}
	if req.OperationType != "" {
		filters["operationType"] = req.OperationType
	}
	if req.OperationObject != "" {
		filters["operationObject"] = req.OperationObject
	}
	if req.Method != "" {
		filters["method"] = req.Method
	}
	if req.StartTime != "" {
		filters["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		filters["endTime"] = req.EndTime
	}

	params := task.ExcelExportParams{
		TaskParamsReq: task.TaskParamsReq{
			Module: consts.TaskModuleOperationLog,
		},
		Filters: filters,
	}

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "构造导出参数失败", err)
	}

	// 创建任务记录
	now := time.Now().Unix()
	taskModel := &taskmodel.AdminTask{
		Name:          fmt.Sprintf("操作日志导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		Type:          consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Status:        consts.TaskStatusPending,
		Params:        sql.NullString{String: string(paramsBytes), Valid: true},
		UserId:        user.UserID,
		ScheduledAt:   0,
		StartedAt:     0,
		FinishedAt:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     0,
	}

	// TODO(phase2-task-rpc): 跨域写入 Task 域（发起导出任务），Phase 2 拆分后改为调用 task-rpc.CreateTask
	taskId, err := l.svcCtx.Domain.Task.Task.Create(l.ctx, taskModel)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建导出任务失败", err)
	}

	logx.Infof("操作日志导出任务已创建: taskId=%d, userId=%d", taskId, user.UserID)

	// 当前接口不直接返回下载URL，URL 由异步任务执行完成后写入任务结果JSON，由前端在任务列表中查看
	return &types.OperationLogExportResp{
		Url: "",
	}, nil
}
