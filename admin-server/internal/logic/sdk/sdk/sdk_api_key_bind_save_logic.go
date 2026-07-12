// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/sdkclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyBindSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyBindSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyBindSaveLogic {
	return &SdkApiKeyBindSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SdkApiKeyBindSave 薄胶水：sdkInterfaceId==0 过滤、事务包裹的"软删旧绑定+插入新绑定"
// 已经搬进 services/sdk/internal/logic/sdkapikeybindsavelogic.go（调
// SDKService.SaveApiKeyBindings），这里原样透传全部绑定，由 sdk-rpc 侧过滤。
func (l *SdkApiKeyBindSaveLogic) SdkApiKeyBindSave(req *types.SdkApiKeyBindSaveReq) error {
	if req == nil || req.SdkKeyId == 0 {
		return errs.New(errs.CodeBadRequest, "sdkKeyId 不能为空")
	}

	bindings := make([]*sdkclient.SdkApiKeyBinding, 0, len(req.Bindings))
	for _, b := range req.Bindings {
		bindings = append(bindings, &sdkclient.SdkApiKeyBinding{
			SdkInterfaceId:  b.SdkInterfaceId,
			CustomRateLimit: b.CustomRateLimit,
		})
	}

	if _, err := l.svcCtx.SdkRPC.SdkApiKeyBindSave(l.ctx, &sdkclient.SdkApiKeyBindSaveRequest{
		SdkKeyId: req.SdkKeyId,
		Bindings: bindings,
	}); err != nil {
		return errs.WrapGRPCError("保存授权失败", err)
	}

	return nil
}
