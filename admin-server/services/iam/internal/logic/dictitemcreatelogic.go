package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictItemCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemCreateLogic {
	return &DictItemCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictItemCreateLogic) DictItemCreate(in *iam.DictItemCreateRequest) (*iam.Empty, error) {
	if in == nil || in.TypeId == 0 || in.Label == "" || in.Value == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型ID、标签和值不能为空"))
	}

	if _, err := l.svcCtx.Domain.System.DictType.FindByID(l.ctx, in.TypeId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "字典类型不存在", err))
	}

	status := in.Status
	if status == 0 {
		status = 1
	}

	dictItem := systemmodel.AdminDictItem{
		TypeId: in.TypeId,
		Label:  in.Label,
		Value:  in.Value,
		Sort:   in.Sort,
		Status: status,
		Remark: sql.NullString{String: in.Remark, Valid: in.Remark != ""},
	}

	if err := l.svcCtx.Domain.System.DictItem.Create(l.ctx, &dictItem); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建字典项失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	typeID := in.TypeId
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
