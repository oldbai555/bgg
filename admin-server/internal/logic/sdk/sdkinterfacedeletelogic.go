// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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
	repo := repository.NewSdkAdminRepository(l.svcCtx.Repository)
	if err := repo.DeleteInterface(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除接口失败", err)
	}

	return nil
}
