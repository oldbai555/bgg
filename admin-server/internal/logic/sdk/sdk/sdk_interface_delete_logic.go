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

type SdkInterfaceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceDeleteLogic {
	return &SdkInterfaceDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkInterfaceDeleteLogic) SdkInterfaceDelete(req *types.SdkInterfaceDeleteReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "ID 不能为空")
	}
	if _, err := l.svcCtx.SdkRPC.SdkInterfaceDelete(l.ctx, &sdkclient.SdkInterfaceDeleteRequest{Id: req.Id}); err != nil {
		return errs.WrapGRPCError("删除接口失败", err)
	}

	return nil
}
