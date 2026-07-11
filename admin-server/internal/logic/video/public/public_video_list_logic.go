// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicVideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoListLogic {
	return &PublicVideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicVideoListLogic) PublicVideoList(req *types.PublicVideoListReq) (resp *types.PublicVideoListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 设置默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 查询视频列表（只查询采集视频，type=2）
	keyword := req.Content
	list, total, err := l.svcCtx.Domain.Video.Video.FindPage(l.ctx, page, size, keyword, 2)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询视频列表失败", err)
	}

	// 转换为公开视频项（字段较少）
	items := make([]types.PublicVideoItem, 0, len(list))
	for _, v := range list {
		item := types.PublicVideoItem{
			Id:   v.Id,
			Name: v.Name,
		}

		if v.Uuid.Valid {
			item.Uuid = v.Uuid.String
		}
		if v.GodNum.Valid {
			item.GodNum = v.GodNum.String
		}

		items = append(items, item)
	}

	return &types.PublicVideoListResp{
		List:  items,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}
