// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
)

type ApiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiDeleteLogic {
	return &ApiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiDeleteLogic) ApiDelete(req *types.ApiDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "接口ID不能为空")
	}

	apiRepo := iamrepo.NewApiRepository(l.svcCtx.Repository)
	if err := apiRepo.DeleteByID(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除接口失败", err)
	}
	return nil
}
