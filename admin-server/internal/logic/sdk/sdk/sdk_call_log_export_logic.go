// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

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

type SdkCallLogExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkCallLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkCallLogExportLogic {
	return &SdkCallLogExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkCallLogExportLogic) SdkCallLogExport(req *types.SdkCallLogExportReq) (*types.Response, error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	filters := make(map[string]interface{})
	if req.SdkKeyId > 0 {
		filters[consts.TaskFilterSdkKeyId] = req.SdkKeyId
	}
	if req.ApiCode != "" {
		filters[consts.TaskFilterApiCode] = req.ApiCode
	}
	if req.RespCode != 0 {
		filters[consts.TaskFilterRespCode] = req.RespCode
	}
	if req.Ip != "" {
		filters[consts.TaskFilterIP] = req.Ip
	}
	if req.StartTime > 0 {
		filters[consts.TaskFilterStartTime] = req.StartTime
	}
	if req.EndTime > 0 {
		filters[consts.TaskFilterEndTime] = req.EndTime
	}

	paramsJSON, err := json.Marshal(map[string]interface{}{
		"module":  consts.TaskModuleSdkCallLog,
		"filters": filters,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "序列化任务参数失败", err)
	}

	resp, err := l.svcCtx.TaskRPC.SubmitTask(l.ctx, &taskclient.SubmitTaskRequest{
		Name:          fmt.Sprintf("SDK调用日志导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		TaskType:      consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Params:        string(paramsJSON),
		UserId:        user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建导出任务失败", err)
	}

	logx.Infof("SDK调用日志导出任务已创建: taskId=%d, userId=%d", resp.TaskId, user.UserID)

	return &types.Response{
		Code:    0,
		Message: "已创建异步导出任务，请在右下角任务列表查看进度",
	}, nil
}
