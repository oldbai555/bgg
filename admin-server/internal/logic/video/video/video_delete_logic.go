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

type VideoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoDeleteLogic {
	return &VideoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoDelete 薄胶水，实际业务逻辑已经搬进 services/content/internal/logic/videodeletelogic.go。
func (l *VideoDeleteLogic) VideoDelete(req *types.VideoDeleteReq) error {
	_, err := l.svcCtx.ContentRPC.VideoDelete(l.ctx, &contentclient.VideoDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除视频失败", err)
	}
	return nil
}
