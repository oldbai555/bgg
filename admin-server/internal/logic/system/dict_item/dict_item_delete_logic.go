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

type DictItemDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictItemDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemDeleteLogic {
	return &DictItemDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictItemDeleteLogic) DictItemDelete(req *types.DictItemDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "字典项ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictItemDelete(l.ctx, &iamclient.DictItemDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除字典项失败", err)
	}
	return nil
}
