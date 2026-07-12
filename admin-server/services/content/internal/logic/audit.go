package logic

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/content/internal/svc"
)

// recordAuditLog 异步回调 IamCallback.RecordAuditLog 写一条审计日志，语义迁移自单体
// pkg/audit.RecordAuditLog（该函数原本直接持有 gateway 的 *svc.ServiceContext、写的
// admin_audit_log 表物理属于 iam，content-rpc 拆分后两者都拿不到）。原实现本来就是
// "异步、失败只记日志"（内部 `go func` + recover），这里保留同样的语义：RPC 调用放进
// 一个带 recover 的 goroutine，不阻塞、不影响主流程。IP/UA 字段不透传——content-rpc
// 唯一的两处调用点（BlogArticleAudit/BlogArticleAuditUnpublish）原来就是传一个空的
// *http.Request，IP/UA 本来就恒为空。
func recordAuditLog(svcCtx *svc.ServiceContext, userID uint64, username, auditType, auditObject string, detail interface{}) {
	detailJSON := ""
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			detailJSON = string(b)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("记录审计日志时发生 panic: %v, userId=%d, username=%s", r, userID, username)
			}
		}()

		_, err := svcCtx.IamCallback.RecordAuditLog(context.Background(), &iamcallbackpb.RecordAuditLogRequest{
			UserId:      userID,
			Username:    username,
			AuditType:   auditType,
			AuditObject: auditObject,
			DetailJson:  detailJSON,
		})
		if err != nil {
			logx.Errorf("记录审计日志失败: userId=%d, username=%s, auditType=%s, auditObject=%s, error: %v",
				userID, username, auditType, auditObject, err)
		}
	}()
}
