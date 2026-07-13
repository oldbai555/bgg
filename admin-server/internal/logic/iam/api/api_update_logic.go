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

type ApiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdateLogic {
	return &ApiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUpdateLogic) ApiUpdate(req *types.ApiUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "接口ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.ApiUpdate(l.ctx, &iamclient.ApiUpdateRequest{
		Id:          req.Id,
		Name:        req.Name,
		Method:      req.Method,
		Path:        req.Path,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新接口失败", err)
	}
	return nil
}
