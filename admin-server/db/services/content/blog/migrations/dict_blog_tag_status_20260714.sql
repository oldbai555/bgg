-- ============================================
-- 字典SQL增量脚本
-- 模块：博客标签管理
-- 创建时间：2026-07-14
-- 说明：补齐 blog_tag.status 字段建表时就声明、但一直没有建的 blog_tag_status 字典
--       （见 db/services/content/blog/create_table_blog.sql 里 blog_tag.status 列注释）
-- ============================================

INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('博客标签状态', 'blog_tag_status', '博客标签启用/禁用状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_tag_status_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_tag_status' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_tag_status_type_id, '启用', '1', 1, 1, '标签启用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_tag_status_type_id, '禁用', '2', 2, 1, '标签禁用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
