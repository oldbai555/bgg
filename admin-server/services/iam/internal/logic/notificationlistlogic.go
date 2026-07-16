package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationListLogic) NotificationList(in *iam.NotificationListRequest) (*iam.NotificationListResponse, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
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

	list, total, err := l.svcCtx.Domain.System.Notification.FindPage(l.ctx, page, pageSize, in.UserId, in.SourceType, in.ReadStatus)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询消息通知列表失败", err))
	}

	items := make([]*iam.NotificationItem, 0, len(list))
	for _, n := range list {
		items = append(items, &iam.NotificationItem{
			Id:         n.Id,
			UserId:     n.UserId,
			SourceType: n.SourceType,
			SourceId:   n.SourceId,
			Title:      n.Title,
			Content:    n.Content,
			ReadStatus: n.ReadStatus,
			ReadAt:     n.ReadAt,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		})
	}

	return &iam.NotificationListResponse{Total: total, List: items}, nil
}
