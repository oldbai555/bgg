// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video_collect

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/svc"
)

type VideoCollectOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCollectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCollectOptionsLogic {
	return &VideoCollectOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCollectOptionsLogic) VideoCollectOptions() error {
	// todo: add your logic here and delete this line

	return nil
}
