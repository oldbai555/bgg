package logic

import (
	"context"
	"errors"
	"time"

	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/pkg/useragent"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/model/monitoring"
	"postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Auth
func (l *LoginLogic) Login(in *iam.LoginRequest) (*iam.TokenPair, error) {
	if in == nil || in.Username == "" || in.Password == "" {
		l.recordLoginLog(0, in.GetUsername(), in.GetClientIp(), in.GetUserAgent(), "用户名和密码不能为空", false)
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户名和密码不能为空"))
	}

	user, err := l.svcCtx.Domain.IAM.User.FindByUsername(l.ctx, in.Username)
	if err != nil {
		if errors.Is(errors.Unwrap(err), sqlx.ErrNotFound) || errors.Is(err, sqlx.ErrNotFound) {
			l.recordLoginLog(0, in.Username, in.ClientIp, in.UserAgent, "用户名或密码错误", false)
			return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "用户名或密码错误"))
		}
		l.recordLoginLog(0, in.Username, in.ClientIp, in.UserAgent, "查询用户失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户失败", err))
	}

	if user.Status != 1 {
		l.recordLoginLog(user.Id, user.Username, in.ClientIp, in.UserAgent, "账号已被禁用", false)
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "账号已被禁用"))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		l.recordLoginLog(user.Id, user.Username, in.ClientIp, in.UserAgent, "用户名或密码错误", false)
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "用户名或密码错误"))
	}

	accessToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.AccessExpire,
		user.Id,
		user.Username,
		false,
	)
	if err != nil {
		l.recordLoginLog(user.Id, user.Username, in.ClientIp, in.UserAgent, "生成访问令牌失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成访问令牌失败", err))
	}

	refreshToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.RefreshExpire,
		user.Id,
		user.Username,
		true,
	)
	if err != nil {
		l.recordLoginLog(user.Id, user.Username, in.ClientIp, in.UserAgent, "生成刷新令牌失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成刷新令牌失败", err))
	}

	l.recordLoginLog(user.Id, user.Username, in.ClientIp, in.UserAgent, "登录成功", true)

	go l.createUnreadNoticeNotifications(user.Id)

	return &iam.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// recordLoginLog 记录登录日志（异步）。gateway 已经不持有 *http.Request，IP/UA 由
// LoginRequest 显式携带（gateway 侧从真实请求里取出后传入），语义与拆分前一致。
func (l *LoginLogic) recordLoginLog(userId uint64, username, clientIP, userAgentStr string, message string, success bool) {
	browser, os := useragent.ParseUserAgent(userAgentStr)

	status := int64(2)
	if success {
		status = 1
	}

	now := time.Now().Unix()
	loginLog := &monitoring.AdminLoginLog{
		UserId:    userId,
		Username:  username,
		IpAddress: clientIP,
		Location:  "",
		Browser:   browser,
		Os:        os,
		UserAgent: userAgentStr,
		Status:    status,
		Message:   message,
		LoginAt:   now,
		LogoutAt:  0,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Errorf("记录登录日志时发生 panic: %v, userId=%d, username=%s", r, userId, username)
			}
		}()

		if err := l.svcCtx.Domain.Monitoring.LoginLog.Create(context.Background(), loginLog); err != nil {
			l.Errorf("记录登录日志失败: userId=%d, username=%s, status=%d, message=%s, error: %v", userId, username, status, message, err)
		}
	}()
}

// createUnreadNoticeNotifications 为新用户创建未读公告通知
func (l *LoginLogic) createUnreadNoticeNotifications(userID uint64) {
	defer func() {
		if r := recover(); r != nil {
			l.Errorf("创建未读公告通知时发生 panic: %v, userId=%d", r, userID)
		}
	}()

	notices, err := l.svcCtx.Domain.System.Notice.FindPublishedNotReadByUser(context.Background(), userID)
	if err != nil {
		l.Errorf("查询未读公告失败: userId=%d, error: %v", userID, err)
		return
	}
	if len(notices) == 0 {
		return
	}

	now := time.Now().Unix()
	for _, notice := range notices {
		notifications, _, err := l.svcCtx.Domain.System.Notification.FindPage(context.Background(), 1, 100, userID, "notice", -1)
		if err == nil {
			hasNotification := false
			for _, notif := range notifications {
				if notif.SourceId == notice.Id && notif.SourceType == "notice" && notif.DeletedAt == 0 {
					hasNotification = true
					break
				}
			}
			if hasNotification {
				continue
			}
		}

		notification := &system.AdminNotification{
			UserId:     userID,
			SourceType: "notice",
			SourceId:   notice.Id,
			Title:      notice.Title,
			Content:    notice.Content,
			ReadStatus: 1,
			ReadAt:     0,
			CreatedAt:  now,
			UpdatedAt:  now,
			DeletedAt:  0,
		}

		if err := l.svcCtx.Domain.System.Notification.Create(context.Background(), notification); err != nil {
			l.Errorf("创建公告通知失败: userId=%d, noticeId=%d, error: %v", userID, notice.Id, err)
		}
	}
}
