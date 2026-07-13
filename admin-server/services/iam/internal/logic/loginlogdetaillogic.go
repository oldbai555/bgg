package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogDetailLogic {
	return &LoginLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogDetailLogic) LoginLogDetail(in *iam.LoginLogDetailRequest) (*iam.LoginLogItem, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "登录日志ID不能为空"))
	}

	log, err := l.svcCtx.Domain.Monitoring.LoginLog.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询登录日志详情失败", err))
	}

	return &iam.LoginLogItem{
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
	}, nil
}
