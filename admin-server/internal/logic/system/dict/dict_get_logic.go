// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
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

	return l.getDictInternal(req.Code)
}

// PublicDictGet 公共字典查询：仅允许白名单 code（白名单校验是 gateway 自己的权限
// 关切，不下沉到 iam-rpc；底层查询逻辑和 DictGet 完全一样，共用同一个 IamRPC.DictGet）。
func (l *DictGetLogic) PublicDictGet(req *types.DictGetReq) (resp *types.DictGetResp, err error) {
	if req == nil || req.Code == "" {
		return nil, errs.New(errs.CodeBadRequest, "字典类型编码不能为空")
	}

	// 白名单校验：当前仅允许 video_proxy_url，后续可按需扩展
	switch req.Code {
	case consts.DictCodeVideoProxyURL:
		// 允许
	default:
		return nil, errs.New(errs.CodeForbidden, "不支持的字典类型")
	}

	return l.getDictInternal(req.Code)
}

func (l *DictGetLogic) getDictInternal(code string) (resp *types.DictGetResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.DictGet(l.ctx, &iamclient.DictGetRequest{Code: code})
	if err != nil {
		return nil, errs.WrapGRPCError("查询字典失败", err)
	}

	items := make([]types.DictItemItem, 0, len(rpcResp.Items))
	for _, di := range rpcResp.Items {
		items = append(items, types.DictItemItem{
			Id:        di.Id,
			TypeId:    di.TypeId,
			Label:     di.Label,
			Value:     di.Value,
			Sort:      di.Sort,
			Status:    di.Status,
			Remark:    di.Remark,
			CreatedAt: di.CreatedAt,
		})
	}

	return &types.DictGetResp{Code: rpcResp.Code, Items: items}, nil
}
