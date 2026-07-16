// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeListLogic {
	return &NoticeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeListLogic) NoticeList(req *types.NoticeListReq) (resp *types.NoticeListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 10, 100)

	noticeType := req.NoticeType
	if noticeType < 0 {
		noticeType = 0
	}
	status := req.Status
	if status < 0 {
		status = -1
	}

	rpcResp, err := l.svcCtx.IamRPC.NoticeList(l.ctx, &iamclient.NoticeListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Title:      req.Title,
		NoticeType: noticeType,
		Status:     status,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询公告列表失败", err)
	}

	items := make([]types.NoticeItem, 0, len(rpcResp.List))
	for _, n := range rpcResp.List {
		items = append(items, types.NoticeItem{
			Id:          n.Id,
			Title:       n.Title,
			Content:     n.Content,
			NoticeType:  n.NoticeType,
			Status:      n.Status,
			PublishTime: n.PublishTime,
			CreatedBy:   n.CreatedBy,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		})
	}

	return &types.NoticeListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
