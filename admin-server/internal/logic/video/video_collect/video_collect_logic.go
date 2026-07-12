// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package video_collect

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCollectLogic {
	return &VideoCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoCollect 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/videocollectlogic.go。
// 原响应里的 Code/Msg 固定成功文案由 gateway 自己拼装，不跨 RPC 边界传递。
func (l *VideoCollectLogic) VideoCollect(req *types.VideoCollectReq) (resp *types.VideoCollectResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.VideoCollect(l.ctx, &contentclient.VideoCollectRequest{
		Uuid:      req.Uuid,
		PlayerUrl: req.PlayerUrl,
		Name:      req.Name,
		GodNum:    req.GodNum,
		XlzzUrls:  req.XlzzUrls,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("新增失败", err)
	}

	data := types.VideoItem{}
	if rpcResp.Data != nil {
		data = types.VideoItem{
			Id:          rpcResp.Data.Id,
			Uuid:        rpcResp.Data.Uuid,
			Name:        rpcResp.Data.Name,
			Cover:       rpcResp.Data.Cover,
			GodNum:      rpcResp.Data.GodNum,
			Duration:    rpcResp.Data.Duration,
			PlayUrl:     rpcResp.Data.PlayUrl,
			XlzzUrls:    rpcResp.Data.XlzzUrls,
			Description: rpcResp.Data.Description,
			SourceType:  rpcResp.Data.SourceType,
			CreatedAt:   rpcResp.Data.CreatedAt,
			UpdatedAt:   rpcResp.Data.UpdatedAt,
		}
	}

	return &types.VideoCollectResp{Code: 200, Msg: "新增成功", Data: data}, nil
}
