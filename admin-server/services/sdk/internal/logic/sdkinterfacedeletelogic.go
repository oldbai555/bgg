package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkInterfaceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceDeleteLogic {
	return &SdkInterfaceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkInterfaceDeleteLogic) SdkInterfaceDelete(in *sdk.SdkInterfaceDeleteRequest) (*sdk.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "ID 不能为空"))
	}
	if err := l.svcCtx.Admin.DeleteInterface(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除接口失败", err))
	}

	return &sdk.Empty{}, nil
}
