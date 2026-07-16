// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package operation_log

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
//  1. 根据筛选条件构造导出参数（module + filters，和 task-rpc 侧 ExcelExportParams 的
//     JSON 结构对应，见 services/task/internal/domain/task/types.go）
//  2. 调 task-rpc.SubmitTask 创建异步任务记录（task-rpc 拆分后不再直接写 admin_task 表，
//     见 17-async-eventing.md 第 1.4 节"提交路径"）
//  3. 由 task-rpc 的调度器 + GenericExportExecutor 实际生成文件，并在任务结果中写入下载URL
func (l *OperationLogExportLogic) OperationLogExport(req *types.OperationLogExportReq) (*types.OperationLogExportResp, error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	filters := make(map[string]interface{})
	if req.UserId > 0 {
		filters[consts.TaskFilterUserId] = req.UserId
	}
	if req.Username != "" {
		filters[consts.TaskFilterUsername] = req.Username
	}
	if req.OperationType != "" {
		filters[consts.TaskFilterOperationType] = req.OperationType
	}
	if req.OperationObject != "" {
		filters[consts.TaskFilterOperationObj] = req.OperationObject
	}
	if req.Method != "" {
		filters[consts.TaskFilterMethod] = req.Method
	}
	if req.StartTime != "" {
		filters[consts.TaskFilterStartTime] = req.StartTime
	}
	if req.EndTime != "" {
		filters[consts.TaskFilterEndTime] = req.EndTime
	}

	paramsJSON, err := json.Marshal(map[string]interface{}{
		"module":  consts.TaskModuleOperationLog,
		"filters": filters,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "构造导出参数失败", err)
	}

	resp, err := l.svcCtx.TaskRPC.SubmitTask(l.ctx, &taskclient.SubmitTaskRequest{
		Name:          fmt.Sprintf("操作日志导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		TaskType:      consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Params:        string(paramsJSON),
		UserId:        user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建导出任务失败", err)
	}

	logx.Infof("操作日志导出任务已创建: taskId=%d, userId=%d", resp.TaskId, user.UserID)

	// 当前接口不直接返回下载URL，URL 由异步任务执行完成后写入任务结果JSON，由前端在任务列表中查看
	return &types.OperationLogExportResp{
		Url: "",
	}, nil
}
