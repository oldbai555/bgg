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

type VideoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoUpdateLogic {
	return &VideoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoUpdate 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/videoupdatelogic.go。
func (l *VideoUpdateLogic) VideoUpdate(req *types.VideoUpdateReq) error {
	_, err := l.svcCtx.ContentRPC.VideoUpdate(l.ctx, &contentclient.VideoUpdateRequest{
		Id:          req.Id,
		Name:        req.Name,
		Cover:       req.Cover,
		GodNum:      req.GodNum,
		Duration:    req.Duration,
		PlayUrl:     req.PlayUrl,
		XlzzUrls:    req.XlzzUrls,
		Description: req.Description,
	})
	if err != nil {
		return errs.WrapGRPCError("更新视频失败", err)
	}
	return nil
}
