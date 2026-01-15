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

// 审计类型（字典：audit_log_type 的 value）
const (
	AuditTypeBlogArticleAudit     = "blog_article_audit"
	AuditTypeBlogArticleUnpublish = "blog_article_unpublish"
)

// 审计对象（AuditObject）
const (
	AuditObjectBlogArticle = "blog_article"
)

// 通用打点 module 建议值（MetricReportReq.module）
const (
	MetricModuleBlogArticleList   = "blog_article_list"
	MetricModuleBlogArticleDetail = "blog_article_detail"
	MetricModuleVideoList         = "video_list"
	MetricModuleVideoDetail       = "video_detail"
)

// 字典 code 常量
const (
	DictCodeBlogTagNameMaxLength          = "blog_tag_name_max_length"           // 标签名称最大长度
	DictCodeBlogArticleTitleMaxLength     = "blog_article_title_max_length"      // 文章标题最大长度
	DictCodeBlogArticleSummaryLength      = "blog_article_summary_length"        // 文章摘要截断长度
	DictCodeBlogFriendLinkNameMaxLength   = "blog_friend_link_name_max_length"   // 友情链接名称最大长度
	DictCodeBlogFriendLinkUrlMaxLength    = "blog_friend_link_url_max_length"    // 友情链接URL最大长度
	DictCodeBlogFriendLinkRemarkMaxLength = "blog_friend_link_remark_max_length" // 友情链接备注最大长度
	DictCodeBlogSocialInfoNameMaxLength   = "blog_social_info_name_max_length"   // 社交信息名称最大长度
	DictCodeBlogSocialInfoUrlMaxLength    = "blog_social_info_url_max_length"    // 社交信息URL最大长度
	DictCodeBlogSocialInfoRemarkMaxLength = "blog_social_info_remark_max_length" // 社交信息备注最大长度
	DictCodeBlogArticleTopMaxCount        = "blog_article_top_max_count"         // 文章置顶最大数量
)
