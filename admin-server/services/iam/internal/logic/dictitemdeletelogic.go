package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictItemDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemDeleteLogic {
	return &DictItemDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictItemDeleteLogic) DictItemDelete(in *iam.DictItemDeleteRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典项ID不能为空"))
	}

	dictItem, err := l.svcCtx.Domain.System.DictItem.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典项失败", err))
	}

	if err := l.svcCtx.Domain.System.DictItem.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除字典项失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	typeID := dictItem.TypeId
	go func() {
		dictType, err := l.svcCtx.Domain.System.DictType.FindByID(context.Background(), typeID)
		if err == nil {
			if err := cache.DeleteDictItems(context.Background(), dictType.Code); err != nil {
				l.Errorf("清除字典项缓存失败: code=%s, error=%v", dictType.Code, err)
			}
			if err := cache.DeleteDictItemsByType(context.Background(), typeID); err != nil {
				l.Errorf("清除字典项缓存失败: typeId=%d, error=%v", typeID, err)
			}
		}
	}()

	return &iam.Empty{}, nil
}
