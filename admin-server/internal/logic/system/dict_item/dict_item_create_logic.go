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

type DictItemCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictItemCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemCreateLogic {
	return &DictItemCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictItemCreateLogic) DictItemCreate(req *types.DictItemCreateReq) error {
	if req == nil || req.TypeId == 0 || req.Label == "" || req.Value == "" {
		return errs.New(errs.CodeBadRequest, "字典类型ID、标签和值不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictItemCreate(l.ctx, &iamclient.DictItemCreateRequest{
		TypeId: req.TypeId,
		Label:  req.Label,
		Value:  req.Value,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	})
	if err != nil {
		return errs.WrapGRPCError("创建字典项失败", err)
	}
	return nil
}
