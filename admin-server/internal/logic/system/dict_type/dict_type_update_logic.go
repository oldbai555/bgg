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

type DictTypeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictTypeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeUpdateLogic {
	return &DictTypeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeUpdateLogic) DictTypeUpdate(req *types.DictTypeUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "字典类型ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictTypeUpdate(l.ctx, &iamclient.DictTypeUpdateRequest{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新字典类型失败", err)
	}
	return nil
}
