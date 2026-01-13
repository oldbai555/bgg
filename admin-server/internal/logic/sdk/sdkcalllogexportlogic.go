// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"postapocgame/admin-server/internal/task"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
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

	// 构造任务参数
	params := task.ExcelExportParams{
		TaskParamsReq: task.TaskParamsReq{Module: consts.TaskModuleSdkCallLog},
		Filters:       make(map[string]interface{}),
	}

	if req.SdkKeyId > 0 {
		params.Filters[consts.TaskFilterSdkKeyId] = req.SdkKeyId
	}
	if req.ApiCode != "" {
		params.Filters[consts.TaskFilterApiCode] = req.ApiCode
	}
	if req.RespCode != 0 {
		params.Filters[consts.TaskFilterRespCode] = req.RespCode
	}
	if req.Ip != "" {
		params.Filters[consts.TaskFilterIP] = req.Ip
	}
	if req.StartTime > 0 {
		params.Filters[consts.TaskFilterStartTime] = req.StartTime
	}
	if req.EndTime > 0 {
		params.Filters[consts.TaskFilterEndTime] = req.EndTime
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "序列化任务参数失败", err)
	}

	now := time.Now().Unix()
	taskModel := &model.AdminTask{
		Name:          fmt.Sprintf("SDK调用日志导出_%s", time.Now().Format("2006-01-02 15:04:05")),
		Type:          consts.TaskTypeExcelExport,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Status:        consts.TaskStatusPending,
		Params:        sql.NullString{String: string(paramsJSON), Valid: true},
		UserId:        user.UserID,
		ScheduledAt:   0,
		StartedAt:     0,
		FinishedAt:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     0,
	}

	taskRepo := repository.NewTaskRepository(l.svcCtx.Repository)
	_, err = taskRepo.Create(l.ctx, taskModel)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建导出任务失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "已创建异步导出任务，请在右下角任务列表查看进度",
	}, nil
}
