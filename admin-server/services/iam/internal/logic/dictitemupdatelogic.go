package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictItemUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemUpdateLogic {
	return &DictItemUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictItemUpdateLogic) DictItemUpdate(in *iam.DictItemUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典项ID不能为空"))
	}

	dictItem, err := l.svcCtx.Domain.System.DictItem.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典项失败", err))
	}

	if in.Label != "" {
		dictItem.Label = in.Label
	}
	if in.Value != "" {
		dictItem.Value = in.Value
	}
	if in.Status == 0 || in.Status == 1 {
		dictItem.Status = in.Status
	}
	if in.Sort >= 0 {
		dictItem.Sort = in.Sort
	}

	if err := l.svcCtx.Domain.System.DictItem.Update(l.ctx, dictItem); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新字典项失败", err))
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
