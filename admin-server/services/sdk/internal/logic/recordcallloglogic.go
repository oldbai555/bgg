package logic

import (
	"context"
	"database/sql"

	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordCallLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecordCallLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordCallLogLogic {
	return &RecordCallLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecordCallLog 从 internal/middleware/sdkcalllogmiddleware.go 的 Handle 方法体搬迁而来。
// 原代码用 `_ = logRepo.SaveCallLog(...)` 丢弃错误——这里如实把错误返回给 gRPC 调用方，
// "失败不影响本次 SDK 调用"这个既有语义由调用方（gateway 侧 SDKCallLogMiddleware）负责
// 丢弃/记日志，不在 callee 这一层悄悄吞掉，保持 RPC 契约的诚实。
func (l *RecordCallLogLogic) RecordCallLog(in *sdk.RecordCallLogRequest) (*sdk.Empty, error) {
	err := l.svcCtx.Public.SaveCallLog(l.ctx, &sdkmodel.SdkCallLog{
		SdkKeyId:       in.SdkKeyId,
		SdkInterfaceId: in.SdkInterfaceId,
		ApiCode:        in.ApiCode,
		Path:           in.Path,
		Method:         in.Method,
		Ip:             in.Ip,
		UserAgent:      in.UserAgent,
		ReqBody:        nullString(in.ReqBody),
		RespBody:       nullString(in.RespBody),
		RespCode:       in.RespCode,
		DurationMs:     in.DurationMs,
	})
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &sdk.Empty{}, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
