// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_video

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
)

type PublicVideoDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicVideoDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoDetailLogic {
	return &PublicVideoDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicVideoDetailLogic) PublicVideoDetail(req *types.PublicVideoDetailReq) (resp *types.PublicVideoDetailResp, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "视频ID不能为空")
	}

	// 查询视频详情
	videoRepo := repository.NewVideoRepository(l.svcCtx.Repository)
	video, err := videoRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeNotFound, "视频不存在", err)
	}

	// 只返回采集视频（type=2）
	if video.Type != 2 {
		return nil, errs.New(errs.CodeNotFound, "视频不存在")
	}

	// 转换为响应类型
	resp = &types.PublicVideoDetailResp{
		Id:         video.Id,
		Name:       video.Name,
		PlayUrl:    video.PlayUrl,
		Duration:   video.Duration,
		SourceType: video.Type,
		CreatedAt:  video.CreatedAt,
		UpdatedAt:  video.UpdatedAt,
	}

	// 处理可选字段
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
