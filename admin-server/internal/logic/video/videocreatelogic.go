// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *VideoCreateLogic) VideoCreate(req *types.VideoCreateReq) error {
	if req == nil || req.Name == "" || req.PlayUrl == "" {
		return errs.New(errs.CodeBadRequest, "视频名称和播放链接不能为空")
	}

	now := time.Now().Unix()
	video := model.Video{
		Name:      req.Name,
		Duration:  req.Duration,
		PlayUrl:   req.PlayUrl,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}

	// 处理可选字段
	if req.Cover != "" {
		video.Cover = sql.NullString{String: req.Cover, Valid: true}
	}
	if req.Description != "" {
		video.Description = sql.NullString{String: req.Description, Valid: true}
	}

	repo := repository.NewVideoRepository(l.svcCtx.Repository)
	if err := repo.Create(l.ctx, &video); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建视频失败", err)
	}

	return nil
}
