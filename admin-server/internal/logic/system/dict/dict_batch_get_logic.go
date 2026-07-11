// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictBatchGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictBatchGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictBatchGetLogic {
	return &DictBatchGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictBatchGetLogic) DictBatchGet(req *types.DictBatchGetReq) (resp *types.DictBatchGetResp, err error) {
	if req == nil || len(req.Codes) == 0 {
		return nil, errs.New(errs.CodeBadRequest, "字典类型编码列表不能为空")
	}

	// 初始化返回结果
	result := make(map[string]types.DictGetResp)
	cache := l.svcCtx.Repository.BusinessCache

	// 遍历每个字典类型编码
	for _, code := range req.Codes {
		if code == "" {
			continue
		}

		// 查询字典类型
		dictType, err := l.svcCtx.Domain.System.DictType.FindByCode(l.ctx, code)
		if err != nil {
			// 如果字典类型不存在，记录日志但继续处理其他字典
			l.Errorf("查询字典类型失败: code=%s, error=%v", code, err)
			continue
		}

		// 尝试从缓存获取
		var cachedItems []types.DictItemItem
		err = cache.GetDictItems(l.ctx, code, &cachedItems)
		if err == nil && len(cachedItems) > 0 {
			// 缓存命中，直接使用
			result[code] = types.DictGetResp{
				Code:  dictType.Code,
				Items: cachedItems,
			}
			continue
		}

		// 缓存未命中，从数据库查询
		items, err := l.svcCtx.Domain.System.DictItem.FindByTypeID(l.ctx, dictType.Id)
		if err != nil {
			// 查询失败，记录日志但继续处理其他字典
			l.Errorf("查询字典项失败: code=%s, error=%v", code, err)
			continue
		}

		// 转换为响应类型
		dictItems := make([]types.DictItemItem, 0, len(items))
		for _, di := range items {
			remark := ""
			if di.Remark.Valid {
				remark = di.Remark.String
			}
			dictItems = append(dictItems, types.DictItemItem{
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

		// 添加到结果中
		result[code] = types.DictGetResp{
			Code:  dictType.Code,
			Items: dictItems,
		}

		// 异步写入缓存（不阻塞返回）
		go func(code string, typeId uint64, items []types.DictItemItem) {
			if err := cache.SetDictItems(context.Background(), code, items); err != nil {
				l.Errorf("设置字典项缓存失败: code=%s, error=%v", code, err)
			}
			// 同时按 type_id 缓存
			if err := cache.SetDictItemsByType(context.Background(), typeId, items); err != nil {
				l.Errorf("设置字典项缓存失败: typeId=%d, error=%v", typeId, err)
			}
		}(code, dictType.Id, dictItems)
	}

	return &types.DictBatchGetResp{
		Dicts: result,
	}, nil
}
