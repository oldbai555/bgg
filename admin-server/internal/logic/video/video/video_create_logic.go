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

type VideoCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCreateLogic {
	return &VideoCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoCreate 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/videocreatelogic.go。
func (l *VideoCreateLogic) VideoCreate(req *types.VideoCreateReq) error {
	_, err := l.svcCtx.ContentRPC.VideoCreate(l.ctx, &contentclient.VideoCreateRequest{
		Name:        req.Name,
		Cover:       req.Cover,
		GodNum:      req.GodNum,
		Duration:    req.Duration,
		PlayUrl:     req.PlayUrl,
		XlzzUrls:    req.XlzzUrls,
		Description: req.Description,
		SourceType:  req.SourceType,
	})
	if err != nil {
		return errs.WrapGRPCError("创建视频失败", err)
	}
	return nil
}
