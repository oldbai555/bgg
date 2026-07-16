// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationListLogic) NotificationList(req *types.NotificationListReq) (resp *types.NotificationListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 10, 100)

	rpcResp, err := l.svcCtx.IamRPC.NotificationList(l.ctx, &iamclient.NotificationListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		UserId:     user.UserID,
		SourceType: req.SourceType,
		ReadStatus: req.ReadStatus,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询消息通知列表失败", err)
	}

	items := make([]types.NotificationItem, 0, len(rpcResp.List))
	for _, n := range rpcResp.List {
		items = append(items, types.NotificationItem{
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

	return &types.NotificationListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
