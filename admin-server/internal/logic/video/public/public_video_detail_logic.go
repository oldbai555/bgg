// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicVideoDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicVideoDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoDetailLogic {
	return &PublicVideoDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicVideoDetail 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicvideodetaillogic.go。
func (l *PublicVideoDetailLogic) PublicVideoDetail(req *types.PublicVideoDetailReq) (resp *types.PublicVideoDetailResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicVideoDetail(l.ctx, &contentclient.PublicVideoDetailRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询视频详情失败", err)
	}

	return &types.PublicVideoDetailResp{
		Id:          rpcResp.Id,
		Uuid:        rpcResp.Uuid,
		Name:        rpcResp.Name,
		Cover:       rpcResp.Cover,
		GodNum:      rpcResp.GodNum,
		Duration:    rpcResp.Duration,
		PlayUrl:     rpcResp.PlayUrl,
		XlzzUrls:    rpcResp.XlzzUrls,
		Description: rpcResp.Description,
		SourceType:  rpcResp.SourceType,
		CreatedAt:   rpcResp.CreatedAt,
		UpdatedAt:   rpcResp.UpdatedAt,
	}, nil
}
