// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package video

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoList 薄胶水：统一分页参数（原 logicutil.NormalizePage）的默认值/上限处理已经下沉到
// content-rpc 侧内联实现，实际业务逻辑搬进 services/content/internal/logic/videolistlogic.go。
func (l *VideoListLogic) VideoList(req *types.VideoListReq) (resp *types.VideoListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.VideoList(l.ctx, &contentclient.VideoListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Keyword:    req.Keyword,
		SourceType: req.SourceType,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询视频列表失败", err)
	}

	items := make([]types.VideoItem, 0, len(rpcResp.List))
	for _, v := range rpcResp.List {
		items = append(items, types.VideoItem{
			Id:          v.Id,
			Uuid:        v.Uuid,
			Name:        v.Name,
			Cover:       v.Cover,
			GodNum:      v.GodNum,
			Duration:    v.Duration,
			PlayUrl:     v.PlayUrl,
			XlzzUrls:    v.XlzzUrls,
			Description: v.Description,
			SourceType:  v.SourceType,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &types.VideoListResp{Total: rpcResp.Total, List: items}, nil
}
