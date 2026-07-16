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

type DictTypeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictTypeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeDeleteLogic {
	return &DictTypeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeDeleteLogic) DictTypeDelete(req *types.DictTypeDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "字典类型ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.DictTypeDelete(l.ctx, &iamclient.DictTypeDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除字典类型失败", err)
	}
	return nil
}
