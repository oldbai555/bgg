// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package login_log

import (
	"context"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogListLogic {
	return &LoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogListLogic) LoginLogList(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	loginLogRepo := repository.NewLoginLogRepository(l.svcCtx.Repository)

	// Status 枚举（字典 login_status）：0 = 全部（不筛选）；1 = 成功；2 = 失败
	// DB 中 admin_login_log.status 与枚举值保持一致：1 = 成功，2 = 失败
	l.Infof("查询登录日志: page=%d, pageSize=%d, userId=%d, username=%s, status=%d, startTime=%s, endTime=%s",
		req.Page, req.PageSize, req.UserId, req.Username, req.Status, req.StartTime, req.EndTime)

	list, total, err := loginLogRepo.FindPage(
		l.ctx,
		int64(req.Page),
		int64(req.PageSize),
		req.UserId,
		req.Username,
		req.Status,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询登录日志列表失败", err)
	}

	items := make([]types.LoginLogItem, 0, len(list))
	for _, log := range list {
		items = append(items, types.LoginLogItem{
			Id:        log.Id,
			UserId:    log.UserId,
			Username:  log.Username,
			IpAddress: log.IpAddress,
			Location:  log.Location,
			Browser:   log.Browser,
			Os:        log.Os,
			UserAgent: log.UserAgent,
			Status:    int(log.Status),
			Message:   log.Message,
			LoginAt:   log.LoginAt,
			LogoutAt:  log.LogoutAt,
			CreatedAt: log.CreatedAt,
		})
	}

	return &types.LoginLogListResp{
		List:     items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
