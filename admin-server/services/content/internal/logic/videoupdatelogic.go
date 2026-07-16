package logic

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoUpdateLogic {
	return &VideoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoUpdate 迁移自 internal/logic/video/video/video_update_logic.go。
func (l *VideoUpdateLogic) VideoUpdate(in *content.VideoUpdateRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	video, err := l.svcCtx.Video.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "视频不存在", err))
	}

	if in.Name != "" {
		video.Name = in.Name
	}
	if in.Duration > 0 {
		video.Duration = in.Duration
	}
	if in.PlayUrl != "" {
		video.PlayUrl = in.PlayUrl
	}
	if in.Cover != "" {
		video.Cover = sql.NullString{String: in.Cover, Valid: true}
	} else if video.Cover.Valid {
		video.Cover = sql.NullString{Valid: false}
	}
	if in.Description != "" {
		video.Description = sql.NullString{String: in.Description, Valid: true}
	} else if video.Description.Valid {
		video.Description = sql.NullString{Valid: false}
	}

	video.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Video.Update(l.ctx, video); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新视频失败", err))
	}

	return &content.Empty{}, nil
}
