package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogListLogic {
	return &LoginLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogListLogic) LoginLogList(in *iam.LoginLogListRequest) (*iam.LoginLogListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	page, pageSize := normalizePage(in.Page, in.PageSize, 20, 100)

	list, total, err := l.svcCtx.Domain.Monitoring.LoginLog.FindPage(
		l.ctx, page, pageSize, in.UserId, in.Username, int(in.Status), in.StartTime, in.EndTime,
	)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询登录日志列表失败", err))
	}

	items := make([]*iam.LoginLogItem, 0, len(list))
	for _, log := range list {
		items = append(items, &iam.LoginLogItem{
			Id:        log.Id,
			UserId:    log.UserId,
			Username:  log.Username,
			IpAddress: log.IpAddress,
			Location:  log.Location,
			Browser:   log.Browser,
			Os:        log.Os,
			UserAgent: log.UserAgent,
			Status:    int64(log.Status),
			Message:   log.Message,
			LoginAt:   log.LoginAt,
			LogoutAt:  log.LogoutAt,
			CreatedAt: log.CreatedAt,
		})
	}

	return &iam.LoginLogListResponse{List: items, Total: total}, nil
}
