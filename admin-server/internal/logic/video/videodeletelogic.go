// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	videorepo "postapocgame/admin-server/internal/repository/video"
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

func (l *VideoDeleteLogic) VideoDelete(req *types.VideoDeleteReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	repo := videorepo.NewVideoRepository(l.svcCtx.Repository)
	if err := repo.DeleteByID(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除视频失败", err)
	}

	return nil
}
