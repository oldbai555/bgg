package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyUpdateLogic {
	return &SdkApiKeyUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyUpdateLogic) SdkApiKeyUpdate(in *sdk.SdkApiKeyUpdateRequest) (*sdk.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "ID 不能为空"))
	}

	data, err := l.svcCtx.Admin.FindSdkKey(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "API Key 不存在", err))
	}

	if strings.TrimSpace(in.Name) != "" {
		data.Name = in.Name
	}
	if in.Status == 1 || in.Status == 2 {
		data.Status = in.Status
	}
	if in.ExpireAt >= 0 {
		data.ExpireAt = in.ExpireAt
	}
	if in.IpWhitelist != "" {
		data.IpWhitelist = in.IpWhitelist
	}
	data.Remark = in.Remark

	if err := l.svcCtx.Admin.UpdateSdkKey(l.ctx, data); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新 API Key 失败", err))
	}

	return &sdk.Empty{}, nil
}
