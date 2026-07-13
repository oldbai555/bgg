// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_type

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictTypeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeCreateLogic {
	return &DictTypeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeCreateLogic) DictTypeCreate(req *types.DictTypeCreateReq) error {
	if req == nil || req.Name == "" || req.Code == "" {
		return errs.New(errs.CodeBadRequest, "字典类型名称和编码不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictTypeCreate(l.ctx, &iamclient.DictTypeCreateRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("创建字典类型失败", err)
	}
	return nil
}
