// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/logic/system/dict"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicDictGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicDictGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicDictGetLogic {
	return &PublicDictGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicDictGet 公共字典查询：仅允许白名单 code（如 video_proxy_url）
// 复用 dict 包的 PublicDictGet 逻辑
func (l *PublicDictGetLogic) PublicDictGet(req *types.DictGetReq) (resp *types.DictGetResp, err error) {
	dictLogic := dict.NewDictGetLogic(l.ctx, l.svcCtx)
	return dictLogic.PublicDictGet(req)
}
