package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictBatchGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictBatchGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictBatchGetLogic {
	return &DictBatchGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictBatchGetLogic) DictBatchGet(in *iam.DictBatchGetRequest) (*iam.DictBatchGetResponse, error) {
	if in == nil || len(in.Codes) == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型编码列表不能为空"))
	}

	result := make(map[string]*iam.DictGetResponse)
	cache := l.svcCtx.Repository.BusinessCache

	for _, code := range in.Codes {
		if code == "" {
			continue
		}

		dictType, err := l.svcCtx.Domain.System.DictType.FindByCode(l.ctx, code)
		if err != nil {
			l.Errorf("查询字典类型失败: code=%s, error=%v", code, err)
			continue
		}

		var cachedItems []*iam.DictItemItem
		if err := cache.GetDictItems(l.ctx, code, &cachedItems); err == nil && len(cachedItems) > 0 {
			result[code] = &iam.DictGetResponse{Code: dictType.Code, Items: cachedItems}
			continue
		}

		items, err := l.svcCtx.Domain.System.DictItem.FindByTypeID(l.ctx, dictType.Id)
		if err != nil {
			l.Errorf("查询字典项失败: code=%s, error=%v", code, err)
			continue
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

		result[code] = &iam.DictGetResponse{Code: dictType.Code, Items: dictItems}

		go func(code string, typeId uint64, items []*iam.DictItemItem) {
			if err := cache.SetDictItems(context.Background(), code, items); err != nil {
				l.Errorf("设置字典项缓存失败: code=%s, error=%v", code, err)
			}
			if err := cache.SetDictItemsByType(context.Background(), typeId, items); err != nil {
				l.Errorf("设置字典项缓存失败: typeId=%d, error=%v", typeId, err)
			}
		}(code, dictType.Id, dictItems)
	}

	return &iam.DictBatchGetResponse{Dicts: result}, nil
}
