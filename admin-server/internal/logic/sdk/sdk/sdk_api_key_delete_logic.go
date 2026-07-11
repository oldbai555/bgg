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

type SdkApiKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyDeleteLogic {
	return &SdkApiKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkApiKeyDeleteLogic) SdkApiKeyDelete(req *types.SdkApiKeyDeleteReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "ID 不能为空")
	}
	if err := l.svcCtx.Domain.SDK.Admin.DeleteSdkKey(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除 API Key 失败", err)
	}

	return nil
}
