// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package video_collect

import (
	"context"

	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
