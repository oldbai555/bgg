// Package consts 复制自 internal/consts/blog.go。按 16-rpc-conventions.md 第 6 节的既定
// 策略：直接复制到各服务自己的 internal/consts，不做成共享包（量很小，维护成本可忽略），
// 后续两边如需变更需要各自同步改。字典 code 常量（DictCodeBlog*）不搬过来——content-rpc
// 拆分后这些字典读取全部换成 services/content/etc/content.yaml 的静态 Limits 配置
// （见 18-service-extraction-runbook.md 2.4 节、services/content/internal/config/config.go）。
package consts

// 博客模块：文章状态（字典：blog_article_status）
const (
	BlogArticleStatusDraft        int64 = 1 // 草稿
	BlogArticleStatusPendingAudit int64 = 2 // 待审核
	BlogArticleStatusAuditPassed  int64 = 3 // 审核通过-未上架
	BlogArticleStatusPublished    int64 = 4 // 上架
	BlogArticleStatusUnpublished  int64 = 5 // 下架
)

// 博客模块：审核状态（字典：blog_article_audit_status）
const (
	BlogArticleAuditStatusNotSubmitted int64 = 1 // 未提交
	BlogArticleAuditStatusPending      int64 = 2 // 待审核
	BlogArticleAuditStatusPassed       int64 = 3 // 审核通过
	BlogArticleAuditStatusRejected     int64 = 4 // 审核驳回
)

// 审计类型（字典：audit_log_type 的 value），通过 IamCallback.RecordAuditLog 回调写入。
const (
	AuditTypeBlogArticleAudit     = "blog_article_audit"
	AuditTypeBlogArticleUnpublish = "blog_article_unpublish"
)

// 审计对象（AuditObject）
const (
	AuditObjectBlogArticle = "blog_article"
)
