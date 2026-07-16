package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEffectiveRateLimitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEffectiveRateLimitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEffectiveRateLimitLogic {
	return &GetEffectiveRateLimitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetEffectiveRateLimit 从 internal/middleware/sdkratelimitmiddleware.go 的 Handle 方法体
// 搬迁而来，去掉了 Redis 滑动窗口计数部分（那部分继续留在 gateway，Redis 全服务共享，
// 见 18-service-extraction-runbook.md 2.2 节）。原代码是拿 apiCode 重新查一次接口，这里
// 直接用调用方（VerifyApiKey 阶段已经拿到）传来的 sdk_interface_id 按主键查，等价但少一次
// 按 api_code 的二次查找。
func (l *GetEffectiveRateLimitLogic) GetEffectiveRateLimit(in *sdk.GetEffectiveRateLimitRequest) (*sdk.GetEffectiveRateLimitResponse, error) {
	iface, err := l.svcCtx.Admin.FindInterface(l.ctx, in.SdkInterfaceId)
	if err != nil || iface == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "接口不存在"))
	}

	limit := l.svcCtx.Public.GetDefaultRateLimit(iface.RateLimitDefault, l.svcCtx.RateLimitDefault)

	binding, _ := l.svcCtx.Public.FindKeyApiBinding(l.ctx, in.SdkKeyId, in.SdkInterfaceId)
	if binding != nil && binding.CustomRateLimit > 0 {
		limit = binding.CustomRateLimit
	}

	return &sdk.GetEffectiveRateLimitResponse{Limit: limit}, nil
}
