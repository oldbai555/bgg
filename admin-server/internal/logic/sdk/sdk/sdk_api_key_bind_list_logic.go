// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyBindListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyBindListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyBindListLogic {
	return &SdkApiKeyBindListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkApiKeyBindListLogic) SdkApiKeyBindList(req *types.SdkApiKeyBindListReq) (resp *types.SdkApiKeyBindListResp, err error) {
	if req == nil || req.SdkKeyId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "sdkKeyId 不能为空")
	}

	list, err := l.svcCtx.Domain.SDK.Admin.ListBindings(l.ctx, req.SdkKeyId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询绑定列表失败", err)
	}

	items := make([]types.SdkApiKeyBindItem, 0, len(list))
	for _, v := range list {
		items = append(items, types.SdkApiKeyBindItem{
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

	return &types.SdkApiKeyBindListResp{List: items}, nil
}
