// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoListLogic) VideoList(req *types.VideoListReq) (resp *types.VideoListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	repo := repository.NewVideoRepository(l.svcCtx.Repository)
	list, total, err := repo.FindPage(l.ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询视频列表失败", err)
	}

	items := make([]types.VideoItem, 0, len(list))
	for _, v := range list {
		item := types.VideoItem{
			Id:        v.Id,
			Name:      v.Name,
			Duration:  v.Duration,
			PlayUrl:   v.PlayUrl,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}

		// 处理可选字段
		if v.Cover.Valid {
			item.Cover = v.Cover.String
		}
		if v.Description.Valid {
			item.Description = v.Description.String
		}

		items = append(items, item)
	}

	return &types.VideoListResp{
		Total: total,
		List:  items,
	}, nil
}
