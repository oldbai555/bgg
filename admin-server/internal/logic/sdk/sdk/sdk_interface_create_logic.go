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

type SdkInterfaceCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceCreateLogic {
	return &SdkInterfaceCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SdkInterfaceCreate 薄胶水：apiCode 自动生成、唯一性校验已经搬进
// services/sdk/internal/logic/sdkinterfacecreatelogic.go。
func (l *SdkInterfaceCreateLogic) SdkInterfaceCreate(req *types.SdkInterfaceCreateReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err := l.svcCtx.SdkRPC.SdkInterfaceCreate(l.ctx, &sdkclient.SdkInterfaceCreateRequest{
		Name:             req.Name,
		Path:             req.Path,
		Method:           req.Method,
		RateLimitDefault: req.RateLimitDefault,
		Status:           req.Status,
		Remark:           req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("创建接口失败", err)
	}

	return nil
}
