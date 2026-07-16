// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
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

	_, err := l.svcCtx.IamRPC.ApiDelete(l.ctx, &iamclient.ApiDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除接口失败", err)
	}
	return nil
}
