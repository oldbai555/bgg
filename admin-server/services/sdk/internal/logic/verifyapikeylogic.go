package logic

import (
	"context"
	"strings"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/consts"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyApiKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyApiKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyApiKeyLogic {
	return &VerifyApiKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VerifyApiKey 从 internal/middleware/sdkauthmiddleware.go 的 Handle 方法体原样搬迁而来
// （HTTP 特有的 response.ErrorCtx/context.WithValue 部分改成结构化返回值）。见
// services/sdk/rpc/sdk.proto 的 VerifyApiKeyResponse 注释：valid=false 时用显式
// code/message 字段透传具体失败原因，不用 gRPC status error，保留外部 SDK 调用方
// 依赖的 7 种具体错误文案。
func (l *VerifyApiKeyLogic) VerifyApiKey(in *sdk.VerifyApiKeyRequest) (*sdk.VerifyApiKeyResponse, error) {
	if strings.TrimSpace(in.ApiKey) == "" || strings.TrimSpace(in.ApiSecret) == "" {
		return deny(errs.CodeUnauthorized, "缺少 API Key 或 Secret"), nil
	}

	sdkKey, err := l.svcCtx.Public.FindKeyByApiKey(l.ctx, in.ApiKey)
	if err != nil {
		return deny(errs.CodeUnauthorized, "无效的 API Key"), nil
	}
	if sdkKey.ApiSecret != in.ApiSecret {
		return deny(errs.CodeUnauthorized, "Secret 不匹配"), nil
	}
	if sdkKey.Status != consts.Open {
		return deny(errs.CodeUnauthorized, "API Key 已被禁用"), nil
	}
	if sdkKey.ExpireAt > 0 && time.Now().Unix() > sdkKey.ExpireAt {
		return deny(errs.CodeUnauthorized, "API Key 已过期"), nil
	}

	if !l.svcCtx.Public.IsIPAllowed(sdkKey.IpWhitelist, in.ClientIp) {
		return deny(errs.CodeForbidden, "IP 不在白名单"), nil
	}

	apiCode := l.svcCtx.Public.BuildInterfaceCode(in.Method, in.Path)
	sdkInterface, err := l.svcCtx.Public.FindInterfaceByCode(l.ctx, apiCode)
	if err != nil || sdkInterface == nil || sdkInterface.Status != consts.Open {
		return deny(errs.CodeForbidden, "接口未开通或已禁用"), nil
	}

	return &sdk.VerifyApiKeyResponse{
		Valid:          true,
		SdkKeyId:       sdkKey.Id,
		ApiKey:         sdkKey.ApiKey,
		SdkInterfaceId: sdkInterface.Id,
		ApiCode:        sdkInterface.ApiCode,
	}, nil
}

func deny(code int, message string) *sdk.VerifyApiKeyResponse {
	return &sdk.VerifyApiKeyResponse{
		Valid:   false,
		Code:    int64(code),
		Message: message,
	}
}
