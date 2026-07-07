// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	systemrepo "postapocgame/admin-server/internal/repository/system"
)

type DictGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictGetLogic {
	return &DictGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictGetLogic) DictGet(req *types.DictGetReq) (resp *types.DictGetResp, err error) {
	if req == nil || req.Code == "" {
		return nil, errs.New(errs.CodeBadRequest, "字典类型编码不能为空")
	}

	return l.getDictInternal(req, false)
}

// PublicDictGet 公共字典查询：仅允许白名单 code
func (l *DictGetLogic) PublicDictGet(req *types.DictGetReq) (resp *types.DictGetResp, err error) {
	if req == nil || req.Code == "" {
		return nil, errs.New(errs.CodeBadRequest, "字典类型编码不能为空")
	}

	// 白名单校验：当前仅允许 video_proxy_url，后续可按需扩展
	switch req.Code {
	case "video_proxy_url":
		// 允许
	default:
		return nil, errs.New(errs.CodeForbidden, "不支持的字典类型")
	}

	return l.getDictInternal(req, true)
}

// getDictInternal 复用现有逻辑，isPublic 预留未来差异化处理扩展
func (l *DictGetLogic) getDictInternal(req *types.DictGetReq, isPublic bool) (resp *types.DictGetResp, err error) {
	// 当前公共/私有逻辑一致，预留 isPublic 以便未来在这里做额外限制或审计

	dictTypeRepo := systemrepo.NewDictTypeRepository(l.svcCtx.Repository)
	dictType, err := dictTypeRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询字典类型失败", err)
	}

	// 尝试从缓存获取字典项列表
	cache := l.svcCtx.Repository.BusinessCache
	var cachedItems []types.DictItemItem
	err = cache.GetDictItems(l.ctx, req.Code, &cachedItems)
	if err == nil {
		// 缓存命中，直接返回
		return &types.DictGetResp{
			Code:  dictType.Code,
			Items: cachedItems,
		}, nil
	}

	// 缓存未命中，从数据库查询
	dictItemRepo := systemrepo.NewDictItemRepository(l.svcCtx.Repository)
	items, err := dictItemRepo.FindByTypeID(l.ctx, dictType.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询字典项失败", err)
	}

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

	resp = &types.DictGetResp{
		Code:  dictType.Code,
		Items: dictItems,
	}

	// 写入缓存（异步，不阻塞返回）
	go func() {
		if err := cache.SetDictItems(context.Background(), req.Code, dictItems); err != nil {
			l.Errorf("设置字典项缓存失败: code=%s, error=%v", req.Code, err)
		}
		// 同时按 type_id 缓存
		if err := cache.SetDictItemsByType(context.Background(), dictType.Id, dictItems); err != nil {
			l.Errorf("设置字典项缓存失败: typeId=%d, error=%v", dictType.Id, err)
		}
	}()

	return resp, nil
}
