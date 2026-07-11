// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyUpdateLogic {
	return &SdkApiKeyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkApiKeyUpdateLogic) SdkApiKeyUpdate(req *types.SdkApiKeyUpdateReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "ID 不能为空")
	}

	data, err := l.svcCtx.Domain.SDK.Admin.FindSdkKey(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeBadRequest, "API Key 不存在", err)
	}

	if strings.TrimSpace(req.Name) != "" {
		data.Name = req.Name
	}
	if req.Status == 1 || req.Status == 2 {
		data.Status = req.Status
	}
	if req.ExpireAt >= 0 {
		data.ExpireAt = req.ExpireAt
	}
	if req.IpWhitelist != "" {
		data.IpWhitelist = req.IpWhitelist
	}
	data.Remark = req.Remark

	if err := l.svcCtx.Domain.SDK.Admin.UpdateSdkKey(l.ctx, data); err != nil {
		return errs.Wrap(errs.CodeInternalError, "更新 API Key 失败", err)
	}

	return nil
}
