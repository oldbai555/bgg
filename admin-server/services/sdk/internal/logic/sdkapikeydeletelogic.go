package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyDeleteLogic {
	return &SdkApiKeyDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyDeleteLogic) SdkApiKeyDelete(in *sdk.SdkApiKeyDeleteRequest) (*sdk.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "ID 不能为空"))
	}
	if err := l.svcCtx.Admin.DeleteSdkKey(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除 API Key 失败", err))
	}

	return &sdk.Empty{}, nil
}
