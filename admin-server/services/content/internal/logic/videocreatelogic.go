package logic

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	videomodel "postapocgame/admin-server/services/content/internal/model/video"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCreateLogic {
	return &VideoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoCreate 迁移自 internal/logic/video/video/video_create_logic.go。
func (l *VideoCreateLogic) VideoCreate(in *content.VideoCreateRequest) (*content.Empty, error) {
	if in.Name == "" || in.PlayUrl == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "视频名称和播放链接不能为空"))
	}

	now := time.Now().Unix()
	video := videomodel.Video{
		Name:      in.Name,
		Duration:  in.Duration,
		PlayUrl:   in.PlayUrl,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}
	if in.Cover != "" {
		video.Cover = sql.NullString{String: in.Cover, Valid: true}
	}
	if in.Description != "" {
		video.Description = sql.NullString{String: in.Description, Valid: true}
	}

	if err := l.svcCtx.Video.Create(l.ctx, &video); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建视频失败", err))
	}

	return &content.Empty{}, nil
}
