package logic

import (
	"context"
	"encoding/json"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Video

// VideoList 迁移自 internal/logic/video/video/video_list_logic.go。原实现用
// logicutil.NormalizePage 统一分页参数——那是 gateway 侧的工具函数，RPC 服务不反向依赖
// gateway internal/，这里内联同等的默认值/上限处理（page 默认 1、pageSize 默认 20、
// 上限 100），和 services/task、services/sdk 各 RPC 服务自己内联分页处理的既有先例一致。
func (l *VideoListLogic) VideoList(in *content.VideoListRequest) (*content.VideoListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	list, total, err := l.svcCtx.Video.FindPage(l.ctx, page, pageSize, in.Keyword, in.SourceType)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询视频列表失败", err))
	}

	items := make([]*content.VideoItem, 0, len(list))
	for _, v := range list {
		item := &content.VideoItem{
			Id:         v.Id,
			Name:       v.Name,
			Duration:   v.Duration,
			PlayUrl:    v.PlayUrl,
			SourceType: v.Type,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		}
		if v.Uuid.Valid {
			item.Uuid = v.Uuid.String
		}
		if v.Cover.Valid {
			item.Cover = v.Cover.String
		}
		if v.GodNum.Valid {
			item.GodNum = v.GodNum.String
		}
		if v.Description.Valid {
			item.Description = v.Description.String
		}
		if v.XlzzUrls.Valid {
			var xlzzUrls []string
			if err := json.Unmarshal([]byte(v.XlzzUrls.String), &xlzzUrls); err == nil {
				item.XlzzUrls = xlzzUrls
			}
		}
		items = append(items, item)
	}

	return &content.VideoListResponse{Total: total, List: items}, nil
}
