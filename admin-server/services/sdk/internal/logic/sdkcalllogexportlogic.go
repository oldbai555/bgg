package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkCallLogExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkCallLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkCallLogExportLogic {
	return &SdkCallLogExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SdkCallLogExport 供单体内嵌的 TaskCallback server（internal/rpcserver/taskcallback/
// server.go 的 fetchSdkCallLog）回调，取代原来直接持有 SdkAdminRepository.ExportCallLogs
// 的方式——sdk_call_log 表现在物理上属于 sdk-rpc。
func (l *SdkCallLogExportLogic) SdkCallLogExport(in *sdk.SdkCallLogExportRequest) (*sdk.SdkCallLogExportResponse, error) {
	list, err := l.svcCtx.Admin.ExportCallLogs(l.ctx, in.MaxRows, in.SdkKeyId, in.ApiCode, in.RespCode, in.Ip, in.StartTime, in.EndTime)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询SDK调用日志失败", err))
	}

	items := make([]*sdk.SdkCallLogItem, 0, len(list))
	for _, v := range list {
		items = append(items, &sdk.SdkCallLogItem{
			Id:             v.Id,
			SdkKeyId:       v.SdkKeyId,
			SdkInterfaceId: v.SdkInterfaceId,
			ApiCode:        v.ApiCode,
			Path:           v.Path,
			Method:         v.Method,
			Ip:             v.Ip,
			RespCode:       v.RespCode,
			DurationMs:     v.DurationMs,
			CreatedAt:      v.CreatedAt,
		})
	}

	return &sdk.SdkCallLogExportResponse{List: items}, nil
}
