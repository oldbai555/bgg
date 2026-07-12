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

// SdkApiKeyUpdate 薄胶水：找旧记录、按字段合并更新的业务逻辑已经搬进
// services/sdk/internal/logic/sdkapikeyupdatelogic.go。
func (l *SdkApiKeyUpdateLogic) SdkApiKeyUpdate(req *types.SdkApiKeyUpdateReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err := l.svcCtx.SdkRPC.SdkApiKeyUpdate(l.ctx, &sdkclient.SdkApiKeyUpdateRequest{
		Id:          req.Id,
		Name:        req.Name,
		Status:      req.Status,
		ExpireAt:    req.ExpireAt,
		IpWhitelist: req.IpWhitelist,
		Remark:      req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("更新 API Key 失败", err)
	}

	return nil
}
