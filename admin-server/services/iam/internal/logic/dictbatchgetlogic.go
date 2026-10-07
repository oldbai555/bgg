package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
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

// DictBatchGet 后台登录后一次要拉几十个字典，外部库单次往返就要几十到上百毫秒，
// 逐个编码查库会超过网关调 iam 的 RPC 超时；这里固定最多两次查库：字典类型按编码 IN 一次，
// 缓存未命中的字典项按 type_id IN 一次。不存在的编码直接缺省，与单个查询的语义一致。
func (l *DictBatchGetLogic) DictBatchGet(in *iam.DictBatchGetRequest) (*iam.DictBatchGetResponse, error) {
	if in == nil || len(in.Codes) == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型编码列表不能为空"))
	}

	result := make(map[string]*iam.DictGetResponse)
	codes := uniqueNonEmpty(in.Codes)
	if len(codes) == 0 {
		return &iam.DictBatchGetResponse{Dicts: result}, nil
	}

	dictTypes, err := l.svcCtx.Domain.System.DictType.FindByCodes(l.ctx, codes)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询字典类型失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	missed := make([]systemmodel.AdminDictType, 0, len(dictTypes))
	for _, dt := range dictTypes {
		var cachedItems []*iam.DictItemItem
		if err := cache.GetDictItems(l.ctx, dt.Code, &cachedItems); err == nil && len(cachedItems) > 0 {
			result[dt.Code] = &iam.DictGetResponse{Code: dt.Code, Items: cachedItems}
			continue
		}
		missed = append(missed, dt)
	}
	if len(missed) == 0 {
		return &iam.DictBatchGetResponse{Dicts: result}, nil
	}

	typeIDs := make([]uint64, 0, len(missed))
	for _, dt := range missed {
		typeIDs = append(typeIDs, dt.Id)
	}
	items, err := l.svcCtx.Domain.System.DictItem.FindByTypeIDs(l.ctx, typeIDs)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询字典项失败", err))
	}

	grouped := groupDictItems(items)
	for _, dt := range missed {
		dictItems := grouped[dt.Id]
		if dictItems == nil {
			dictItems = []*iam.DictItemItem{}
		}
		result[dt.Code] = &iam.DictGetResponse{Code: dt.Code, Items: dictItems}

		go func(code string, typeID uint64, items []*iam.DictItemItem) {
			if err := cache.SetDictItems(context.Background(), code, items); err != nil {
				l.Errorf("设置字典项缓存失败: code=%s, error=%v", code, err)
			}
			if err := cache.SetDictItemsByType(context.Background(), typeID, items); err != nil {
				l.Errorf("设置字典项缓存失败: typeId=%d, error=%v", typeID, err)
			}
		}(dt.Code, dt.Id, dictItems)
	}

	return &iam.DictBatchGetResponse{Dicts: result}, nil
}

func uniqueNonEmpty(codes []string) []string {
	seen := make(map[string]struct{}, len(codes))
	out := make([]string, 0, len(codes))
	for _, c := range codes {
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		out = append(out, c)
	}
	return out
}

// groupDictItems 按 type_id 分组，保持查询返回的顺序（sort asc, id asc）
func groupDictItems(items []systemmodel.AdminDictItem) map[uint64][]*iam.DictItemItem {
	grouped := make(map[uint64][]*iam.DictItemItem)
	for _, di := range items {
		remark := ""
		if di.Remark.Valid {
			remark = di.Remark.String
		}
		grouped[di.TypeId] = append(grouped[di.TypeId], &iam.DictItemItem{
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
	return grouped
}
