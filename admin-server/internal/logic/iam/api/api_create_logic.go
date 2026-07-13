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

type ApiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCreateLogic {
	return &ApiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCreateLogic) ApiCreate(req *types.ApiCreateReq) error {
	if req == nil || req.Name == "" || req.Method == "" || req.Path == "" {
		return errs.New(errs.CodeBadRequest, "接口名称、方法和路径不能为空")
	}

	_, err := l.svcCtx.IamRPC.ApiCreate(l.ctx, &iamclient.ApiCreateRequest{
		Name:        req.Name,
		Method:      req.Method,
		Path:        req.Path,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("创建接口失败", err)
	}
	return nil
}
