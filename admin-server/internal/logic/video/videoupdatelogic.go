// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	videorepo "postapocgame/admin-server/internal/repository/video"
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

func (l *VideoUpdateLogic) VideoUpdate(req *types.VideoUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	repo := videorepo.NewVideoRepository(l.svcCtx.Repository)
	video, err := repo.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeNotFound, "视频不存在", err)
	}

	// 更新字段（只更新提供的字段）
	if req.Name != "" {
		video.Name = req.Name
	}
	if req.Duration > 0 {
		video.Duration = req.Duration
	}
	if req.PlayUrl != "" {
		video.PlayUrl = req.PlayUrl
	}
	if req.Cover != "" {
		video.Cover = sql.NullString{String: req.Cover, Valid: true}
	} else if req.Cover == "" && video.Cover.Valid {
		// 如果传入空字符串，表示清空
		video.Cover = sql.NullString{Valid: false}
	}
	if req.Description != "" {
		video.Description = sql.NullString{String: req.Description, Valid: true}
	} else if req.Description == "" && video.Description.Valid {
		video.Description = sql.NullString{Valid: false}
	}

	video.UpdatedAt = time.Now().Unix()

	if err := repo.Update(l.ctx, video); err != nil {
		return errs.Wrap(errs.CodeInternalError, "更新视频失败", err)
	}

	return nil
}
