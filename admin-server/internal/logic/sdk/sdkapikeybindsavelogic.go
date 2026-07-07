// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model/sdk"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
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

func (l *SdkApiKeyBindSaveLogic) SdkApiKeyBindSave(req *types.SdkApiKeyBindSaveReq) error {
	if req == nil || req.SdkKeyId == 0 {
		return errs.New(errs.CodeBadRequest, "sdkKeyId 不能为空")
	}
	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)

	bindings := make([]sdk.SdkKeyApi, 0, len(req.Bindings))
	for _, b := range req.Bindings {
		if b.SdkInterfaceId == 0 {
			continue
		}
		bindings = append(bindings, sdk.SdkKeyApi{
			SdkKeyId:        req.SdkKeyId,
			SdkInterfaceId:  b.SdkInterfaceId,
			CustomRateLimit: b.CustomRateLimit,
			DeletedAt:       0,
		})
	}

	if err := repo.SaveBindings(l.ctx, req.SdkKeyId, bindings); err != nil {
		return errs.Wrap(errs.CodeInternalError, "保存授权失败", err)
	}

	return nil
}
