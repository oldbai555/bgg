package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeListLogic {
	return &NoticeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NoticeListLogic) NoticeList(in *iam.NoticeListRequest) (*iam.NoticeListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	noticeType := in.NoticeType
	if noticeType < 0 {
		noticeType = 0
	}
	status := in.Status
	if status < 0 {
		status = -1
	}

	list, total, err := l.svcCtx.Domain.System.Notice.FindPage(l.ctx, page, pageSize, in.Title, noticeType, status)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询公告列表失败", err))
	}

	items := make([]*iam.NoticeItem, 0, len(list))
	for _, n := range list {
		items = append(items, &iam.NoticeItem{
			Id:          n.Id,
			Title:       n.Title,
			Content:     n.Content,
			NoticeType:  n.Type,
			Status:      n.Status,
			PublishTime: n.PublishTime,
			CreatedBy:   n.CreatedBy,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		})
	}

	return &iam.NoticeListResponse{Total: total, List: items}, nil
}
