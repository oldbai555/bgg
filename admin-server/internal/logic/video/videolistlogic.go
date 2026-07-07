// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"encoding/json"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	videorepo "postapocgame/admin-server/internal/repository/video"
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

	// 统一分页参数
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	repo := videorepo.NewVideoRepository(l.svcCtx.Repository)
	// 支持sourceType筛选（0=全部，1=手动添加，2=采集）
	sourceType := req.SourceType
	list, total, err := repo.FindPage(l.ctx, req.Page, req.PageSize, req.Keyword, sourceType)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询视频列表失败", err)
	}

	items := make([]types.VideoItem, 0, len(list))
	for _, v := range list {
		item := types.VideoItem{
			Id:         v.Id,
			Name:       v.Name,
			Duration:   v.Duration,
			PlayUrl:    v.PlayUrl,
			SourceType: v.Type,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		}

		// 处理可选字段
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

	return &types.VideoListResp{
		Total: total,
		List:  items,
	}, nil
}
