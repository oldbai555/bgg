package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictGetLogic {
	return &DictGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DictGet 后台/公共两个入口共用（公共入口的白名单校验在 gateway 侧
// public_dict_get_logic.go 完成，见 internal/logic/misc/public/public_dict_get_logic.go
// 迁移前的实现），这里只有一份查询逻辑。
func (l *DictGetLogic) DictGet(in *iam.DictGetRequest) (*iam.DictGetResponse, error) {
	if in == nil || in.Code == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型编码不能为空"))
	}

	dictType, err := l.svcCtx.Domain.System.DictType.FindByCode(l.ctx, in.Code)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典类型失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	var cachedItems []*iam.DictItemItem
	if err := cache.GetDictItems(l.ctx, in.Code, &cachedItems); err == nil {
		return &iam.DictGetResponse{Code: dictType.Code, Items: cachedItems}, nil
	}

	items, err := l.svcCtx.Domain.System.DictItem.FindByTypeID(l.ctx, dictType.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典项失败", err))
	}

	dictItems := make([]*iam.DictItemItem, 0, len(items))
	for _, di := range items {
		remark := ""
		if di.Remark.Valid {
			remark = di.Remark.String
		}
		dictItems = append(dictItems, &iam.DictItemItem{
			Id:        di.Id,
			TypeId:    di.TypeId,
			Label:     di.Label,
			Value:     di.Value,
			Sort:      di.Sort,
			Status:    di.Status,
			Remark:    remark,
			CreatedAt: di.CreatedAt,
		})
	}

	resp := &iam.DictGetResponse{Code: dictType.Code, Items: dictItems}

	code, typeID := in.Code, dictType.Id
	go func() {
		if err := cache.SetDictItems(context.Background(), code, dictItems); err != nil {
			l.Errorf("设置字典项缓存失败: code=%s, error=%v", code, err)
		}
		if err := cache.SetDictItemsByType(context.Background(), typeID, dictItems); err != nil {
			l.Errorf("设置字典项缓存失败: typeId=%d, error=%v", typeID, err)
		}
	}()

	return resp, nil
}
