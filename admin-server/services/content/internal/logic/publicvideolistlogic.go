package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoListLogic {
	return &PublicVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公共视频展示

// PublicVideoList 迁移自 internal/logic/video/public/public_video_list_logic.go。
func (l *PublicVideoListLogic) PublicVideoList(in *content.PublicVideoListRequest) (*content.PublicVideoListResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	size := in.Size
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 只查询采集视频（type=2）
	list, total, err := l.svcCtx.Video.FindPage(l.ctx, page, size, in.Content, 2)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询视频列表失败", err))
	}

	items := make([]*content.PublicVideoItem, 0, len(list))
	for _, v := range list {
		item := &content.PublicVideoItem{Id: v.Id, Name: v.Name}
		if v.Uuid.Valid {
			item.Uuid = v.Uuid.String
		}
		if v.GodNum.Valid {
			item.GodNum = v.GodNum.String
		}
		items = append(items, item)
	}

	return &content.PublicVideoListResponse{List: items, Page: page, Size: size, Total: total}, nil
}
