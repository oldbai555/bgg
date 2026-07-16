package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyBindSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyBindSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyBindSaveLogic {
	return &SdkApiKeyBindSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyBindSaveLogic) SdkApiKeyBindSave(in *sdk.SdkApiKeyBindSaveRequest) (*sdk.Empty, error) {
	if in.SdkKeyId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "sdkKeyId 不能为空"))
	}
	bindings := make([]sdkmodel.SdkKeyApi, 0, len(in.Bindings))
	for _, b := range in.Bindings {
		if b.SdkInterfaceId == 0 {
			continue
		}
		bindings = append(bindings, sdkmodel.SdkKeyApi{
			SdkKeyId:        in.SdkKeyId,
			SdkInterfaceId:  b.SdkInterfaceId,
			CustomRateLimit: b.CustomRateLimit,
			DeletedAt:       0,
		})
	}

	if err := l.svcCtx.Service.SaveApiKeyBindings(l.ctx, in.SdkKeyId, bindings); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "保存授权失败", err))
	}

	return &sdk.Empty{}, nil
}
