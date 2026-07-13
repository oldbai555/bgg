package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"postapocgame/admin-server/internal/svc"
	pb "postapocgame/admin-server/pkg/iamcallback/pb"
	jwthelper "postapocgame/admin-server/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuditType 审计类型常量
const (
	AuditTypePermissionAssign = "permission_assign" // 权限分配
	AuditTypeRoleChange       = "role_change"       // 角色变更
	AuditTypeConfigModify     = "config_modify"     // 配置修改
	AuditTypeDataDelete       = "data_delete"       // 数据删除
)

// AuditObject 审计对象常量
const (
	AuditObjectUserRole       = "user_role"       // 用户-角色关联
	AuditObjectRolePermission = "role_permission" // 角色-权限关联
	AuditObjectRole           = "role"            // 角色
	AuditObjectUser           = "user"            // 用户
	AuditObjectPermission     = "permission"      // 权限
	AuditObjectConfig         = "config"          // 配置
)

// RecordAuditLog 记录审计日志（异步）。iam 域拆分成独立服务后，admin_audit_log 表物理上
// 属于 iam-rpc，这里改成回调 iam-rpc 同一进程内注册的 pkg/iamcallback.IamCallback 服务
// （和 content-rpc 的两处调用点走同一个契约），"异步、失败只记日志"的既有语义不变。
// httpReq 可以为 nil，如果为 nil 则 IP 和 User-Agent 为空。
func RecordAuditLog(svcCtx *svc.ServiceContext, ctx context.Context, httpReq *http.Request, auditType, auditObject string, detail interface{}) {
	if svcCtx == nil || ctx == nil {
		logx.Errorf("记录审计日志失败: svcCtx 或 ctx 为 nil")
		return
	}

	// 获取用户信息
	user, ok := jwthelper.FromContext(ctx)
	userId := uint64(0)
	username := ""
	if ok {
		userId = user.UserID
		username = user.Username
	}

	// 获取 IP 地址和 User-Agent
	ip := ""
	userAgent := ""
	if httpReq != nil {
		ip = getClientIP(httpReq)
		userAgent = httpReq.UserAgent()
	}

	// 序列化审计详情
	detailJSON := ""
	if detail != nil {
		if detailBytes, err := json.Marshal(detail); err == nil {
			detailJSON = string(detailBytes)
		}
	}

	// 异步写入日志（使用 goroutine，避免阻塞主流程）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("记录审计日志时发生 panic: %v, userId=%d, username=%s", r, userId, username)
			}
		}()

		_, err := svcCtx.IamCallbackRPC.RecordAuditLog(context.Background(), &pb.RecordAuditLogRequest{
			UserId:      userId,
			Username:    username,
			AuditType:   auditType,
			AuditObject: auditObject,
			DetailJson:  detailJSON,
			IpAddress:   ip,
			UserAgent:   userAgent,
		})
		if err != nil {
			logx.Errorf("记录审计日志失败: userId=%d, username=%s, auditType=%s, auditObject=%s, error: %v",
				userId, username, auditType, auditObject, err)
		} else {
			logx.Infof("成功记录审计日志: userId=%d, username=%s, auditType=%s, auditObject=%s",
				userId, username, auditType, auditObject)
		}
	}()
}

// getClientIP 获取客户端 IP 地址
func getClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	// 优先从 X-Forwarded-For 获取（代理场景）
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 其次从 X-Real-IP 获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 最后从 RemoteAddr 获取
	ip = r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
