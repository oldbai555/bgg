package logic

import (
	"context"
	"encoding/json"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicVideoDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicVideoDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoDetailLogic {
	return &PublicVideoDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicVideoDetail 迁移自 internal/logic/video/public/public_video_detail_logic.go。
func (l *PublicVideoDetailLogic) PublicVideoDetail(in *content.PublicVideoDetailRequest) (*content.PublicVideoDetailResponse, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "视频ID不能为空"))
	}

	video, err := l.svcCtx.Video.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "视频不存在", err))
	}
	// 只返回采集视频（type=2）
	if video.Type != 2 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "视频不存在"))
	}

	resp := &content.PublicVideoDetailResponse{
		Id:         video.Id,
		Name:       video.Name,
		PlayUrl:    video.PlayUrl,
		Duration:   video.Duration,
		SourceType: video.Type,
		CreatedAt:  video.CreatedAt,
		UpdatedAt:  video.UpdatedAt,
	}
	if video.Uuid.Valid {
		resp.Uuid = video.Uuid.String
	}
	if video.Cover.Valid {
		resp.Cover = video.Cover.String
	}
	if video.GodNum.Valid {
		resp.GodNum = video.GodNum.String
	}
	if video.Description.Valid {
		resp.Description = video.Description.String
	}
	if video.XlzzUrls.Valid {
		var xlzzUrls []string
		if err := json.Unmarshal([]byte(video.XlzzUrls.String), &xlzzUrls); err == nil {
			resp.XlzzUrls = xlzzUrls
		}
	}

	return resp, nil
}
