package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyBindListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyBindListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyBindListLogic {
	return &SdkApiKeyBindListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyBindListLogic) SdkApiKeyBindList(in *sdk.SdkApiKeyBindListRequest) (*sdk.SdkApiKeyBindListResponse, error) {
	if in.SdkKeyId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "sdkKeyId 不能为空"))
	}

	list, err := l.svcCtx.Admin.ListBindings(l.ctx, in.SdkKeyId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询绑定列表失败", err))
	}

	items := make([]*sdk.SdkApiKeyBindItem, 0, len(list))
	for _, v := range list {
		items = append(items, &sdk.SdkApiKeyBindItem{
			SdkInterfaceId:  v.SdkInterfaceId,
			ApiCode:         v.ApiCode,
			Name:            v.Name,
			Path:            v.Path,
			Method:          v.Method,
			Bound:           v.Bound,
			RateLimit:       v.RateLimit,
			CustomRateLimit: v.CustomRateLimit,
		})
	}

	return &sdk.SdkApiKeyBindListResponse{List: items}, nil
}
