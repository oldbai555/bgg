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

type SdkApiKeyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyCreateLogic {
	return &SdkApiKeyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SdkApiKeyCreate 薄胶水：Key/Secret 生成与唯一性校验已经搬进
// services/sdk/internal/logic/sdkapikeycreatelogic.go。
func (l *SdkApiKeyCreateLogic) SdkApiKeyCreate(req *types.SdkApiKeyCreateReq) (resp *types.SdkApiKeyCreateResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.SdkRPC.SdkApiKeyCreate(l.ctx, &sdkclient.SdkApiKeyCreateRequest{
		Name:        req.Name,
		Status:      req.Status,
		ExpireAt:    req.ExpireAt,
		IpWhitelist: req.IpWhitelist,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建 API Key 失败", err)
	}

	return &types.SdkApiKeyCreateResp{
		Id:        rpcResp.Id,
		ApiKey:    rpcResp.ApiKey,
		ApiSecret: rpcResp.ApiSecret,
	}, nil
}
