// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_item

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

	// 先查询字典项，获取 typeId
	dictItem, err := l.svcCtx.Domain.System.DictItem.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "查询字典项失败", err)
	}

	if err := l.svcCtx.Domain.System.DictItem.DeleteByID(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除字典项失败", err)
	}

	// 清除字典缓存
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		// 需要获取字典类型的 code 来清除缓存
		dictType, err := l.svcCtx.Domain.System.DictType.FindByID(context.Background(), dictItem.TypeId)
		if err == nil {
			if err := cache.DeleteDictItems(context.Background(), dictType.Code); err != nil {
				l.Errorf("清除字典项缓存失败: code=%s, error=%v", dictType.Code, err)
			}
			if err := cache.DeleteDictItemsByType(context.Background(), dictItem.TypeId); err != nil {
				l.Errorf("清除字典项缓存失败: typeId=%d, error=%v", dictItem.TypeId, err)
			}
		}
	}()

	return nil
}
