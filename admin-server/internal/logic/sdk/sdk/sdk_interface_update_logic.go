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

type SdkInterfaceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceUpdateLogic {
	return &SdkInterfaceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SdkInterfaceUpdate 薄胶水：找旧记录、Path/Method 变更时重新生成 apiCode 并查重的
// 业务逻辑已经搬进 services/sdk/internal/logic/sdkinterfaceupdatelogic.go。
func (l *SdkInterfaceUpdateLogic) SdkInterfaceUpdate(req *types.SdkInterfaceUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "ID 不能为空")
	}

	_, err := l.svcCtx.SdkRPC.SdkInterfaceUpdate(l.ctx, &sdkclient.SdkInterfaceUpdateRequest{
		Id:               req.Id,
		Name:             req.Name,
		Path:             req.Path,
		Method:           req.Method,
		RateLimitDefault: req.RateLimitDefault,
		Status:           req.Status,
		Remark:           req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("更新接口失败", err)
	}

	return nil
}
