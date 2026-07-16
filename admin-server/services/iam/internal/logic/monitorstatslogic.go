package logic

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MonitorStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMonitorStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MonitorStatsLogic {
	return &MonitorStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MonitorStatsLogic) MonitorStats(in *iam.Empty) (*iam.MonitorStatsResponse, error) {
	userCount, err := l.countUsers()
	if err != nil {
		l.Errorf("统计用户数失败: %v", err)
		userCount = 0
	}

	roleCount, err := l.countRoles()
	if err != nil {
		l.Errorf("统计角色数失败: %v", err)
		roleCount = 0
	}

	permissionCount, err := l.countPermissions()
	if err != nil {
		l.Errorf("统计权限数失败: %v", err)
		permissionCount = 0
	}

	menuCount, err := l.countMenus()
	if err != nil {
		l.Errorf("统计菜单数失败: %v", err)
		menuCount = 0
	}

	onlineUserCount := int64(0)
	if onlineResp, err := l.svcCtx.ChatRPC.GetOnlineUserCount(l.ctx, &chatclient.Empty{}); err == nil {
		onlineUserCount = onlineResp.Count
	}

	operationLogCount, err := l.countOperationLogs()
	if err != nil {
		l.Errorf("统计操作日志数失败: %v", err)
		operationLogCount = 0
	}

	loginLogCount, err := l.countLoginLogs()
	if err != nil {
		l.Errorf("统计登录日志数失败: %v", err)
		loginLogCount = 0
	}

	return &iam.MonitorStatsResponse{
		UserCount:         userCount,
		RoleCount:         roleCount,
		PermissionCount:   permissionCount,
		MenuCount:         menuCount,
		OnlineUserCount:   onlineUserCount,
		OperationLogCount: operationLogCount,
		LoginLogCount:     loginLogCount,
	}, nil
}

func (l *MonitorStatsLogic) countUsers() (int64, error) {
	return l.countTable("admin_user")
}

func (l *MonitorStatsLogic) countRoles() (int64, error) {
	return l.countTable("admin_role")
}

func (l *MonitorStatsLogic) countPermissions() (int64, error) {
	return l.countTable("admin_permission")
}

func (l *MonitorStatsLogic) countMenus() (int64, error) {
	return l.countTable("admin_menu")
}

func (l *MonitorStatsLogic) countOperationLogs() (int64, error) {
	return l.countTable("admin_operation_log")
}

func (l *MonitorStatsLogic) countLoginLogs() (int64, error) {
	return l.countTable("admin_login_log")
}

func (l *MonitorStatsLogic) countTable(table string) (int64, error) {
	query, args, err := sq.Select("COUNT(*)").
		From("`" + table + "`").
		Where(sq.Eq{"deleted_at": 0}).
		ToSql()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "统计查询sql生成失败", err)
	}

	var count int64
	if err := l.svcCtx.Repository.DB.QueryRowCtx(l.ctx, &count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}
