// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_item

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictItemUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemUpdateLogic {
	return &DictItemUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictItemUpdateLogic) DictItemUpdate(req *types.DictItemUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "字典项ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictItemUpdate(l.ctx, &iamclient.DictItemUpdateRequest{
		Id:     req.Id,
		Label:  req.Label,
		Value:  req.Value,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("更新字典项失败", err)
	}
	return nil
}
