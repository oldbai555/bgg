-- ============================================
-- 字典SQL增量脚本
-- 模块：博客管理
-- 创建时间：2026-01-14
-- 说明：博客文章状态、审核状态、标签/标题/摘要长度配置及审计类型扩展
-- ============================================

-- 1. 文章业务状态字典：blog_article_status
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客文章状态', 'blog_article_status', '博客文章业务状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_article_status_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_article_status' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_article_status_type_id, '草稿', '1', 1, 1, '草稿状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_status_type_id, '待审核', '2', 2, 1, '已提交待审核', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_status_type_id, '审核通过-未上架', '3', 3, 1, '审核通过但未上架', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_status_type_id, '上架', '4', 4, 1, '已上架，可在公共页面展示', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_status_type_id, '下架', '5', 5, 1, '已下架，不再对外展示', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 文章审核状态字典：blog_article_audit_status
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客文章审核状态', 'blog_article_audit_status', '博客文章审核流程状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_article_audit_status_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_article_audit_status' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_article_audit_status_type_id, '未提交', '1', 1, 1, '尚未提交审核', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_audit_status_type_id, '待审核', '2', 2, 1, '已提交待审核', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_audit_status_type_id, '审核通过', '3', 3, 1, '审核通过', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_audit_status_type_id, '审核驳回', '4', 4, 1, '审核未通过', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 3. 标签名称最大长度配置：blog_tag_name_max_length
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客标签名称长度配置', 'blog_tag_name_max_length', '控制博客标签名称最大长度（字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_tag_name_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_tag_name_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_tag_name_max_len_type_id, '标签名称最大长度', '10', 1, 1, '标签名称最多10个字符（按后端rune长度校验）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 4. 文章标题最大长度配置：blog_article_title_max_length
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客文章标题长度配置', 'blog_article_title_max_length', '控制博客文章标题最大长度（字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_article_title_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_article_title_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_article_title_max_len_type_id, '文章标题最大长度', '100', 1, 1, '文章标题最大长度100个字符（按后端rune长度校验）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 5. 公共展示摘要长度配置：blog_article_summary_length
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客文章摘要展示长度', 'blog_article_summary_length', '控制公共页面文章摘要截断长度（字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_article_summary_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_article_summary_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_article_summary_len_type_id, '文章摘要展示长度', '120', 1, 1, '公共列表/详情摘要默认截断长度120个字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 6. 审计类型扩展：博客文章审核与下架
-- 假设已有审计类型字典 code 为 audit_log_type，如无则后续可按实际 code 调整
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('审计日志类型', 'audit_log_type', '审计日志类型字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @audit_log_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'audit_log_type' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@audit_log_type_id, '博客文章审核', 'blog_article_audit', 10, 1, '文章审核（通过/驳回）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_log_type_id, '博客文章下架', 'blog_article_unpublish', 11, 1, '文章下架（审核员/管理员操作）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

