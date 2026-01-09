// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public

import (
	"context"

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

func (l *PublicDictGetLogic) PublicDictGet(req *types.DictGetReq) (resp *types.DictGetResp, err error) {
	// todo: add your logic here and delete this line

	return
}
