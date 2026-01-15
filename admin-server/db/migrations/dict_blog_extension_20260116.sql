-- ============================================
-- 字典SQL增量脚本
-- 模块：博客扩展功能（友情链接、社交信息、文章置顶）
-- 创建时间：2026-01-16
-- 说明：友情链接、社交信息相关字典配置及文章置顶数量配置
-- ============================================

-- ============================================
-- 1. 友情链接名称最大长度配置：blog_friend_link_name_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('友情链接名称最大长度', 'blog_friend_link_name_max_length', '友情链接名称最大长度（中文字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_friend_link_name_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_friend_link_name_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_friend_link_name_max_len_type_id, '友情链接名称最大长度', '15', 1, 1, '友情链接名称最大15个中文字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 2. 友情链接URL最大长度配置：blog_friend_link_url_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('友情链接URL最大长度', 'blog_friend_link_url_max_length', '友情链接URL最大长度（字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_friend_link_url_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_friend_link_url_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_friend_link_url_max_len_type_id, '友情链接URL最大长度', '255', 1, 1, '友情链接URL最大255个字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 3. 友情链接备注最大长度配置：blog_friend_link_remark_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('友情链接备注最大长度', 'blog_friend_link_remark_max_length', '友情链接备注最大长度（中文字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_friend_link_remark_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_friend_link_remark_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_friend_link_remark_max_len_type_id, '友情链接备注最大长度', '127', 1, 1, '友情链接备注最大127个中文字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 4. 友情链接状态字典：blog_friend_link_status
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('友情链接状态', 'blog_friend_link_status', '友情链接启用/禁用状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_friend_link_status_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_friend_link_status' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_friend_link_status_type_id, '启用', '1', 1, 1, '友情链接启用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_friend_link_status_type_id, '禁用', '2', 2, 1, '友情链接禁用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 5. 社交信息名称最大长度配置：blog_social_info_name_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('社交信息名称最大长度', 'blog_social_info_name_max_length', '社交信息名称最大长度（中文字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_social_info_name_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_social_info_name_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_social_info_name_max_len_type_id, '社交信息名称最大长度', '15', 1, 1, '社交信息名称最大15个中文字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 6. 社交信息URL最大长度配置：blog_social_info_url_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('社交信息URL最大长度', 'blog_social_info_url_max_length', '社交信息URL最大长度（字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_social_info_url_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_social_info_url_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_social_info_url_max_len_type_id, '社交信息URL最大长度', '255', 1, 1, '社交信息URL最大255个字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 7. 社交信息备注最大长度配置：blog_social_info_remark_max_length
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('社交信息备注最大长度', 'blog_social_info_remark_max_length', '社交信息备注最大长度（中文字符数）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_social_info_remark_max_len_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_social_info_remark_max_length' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_social_info_remark_max_len_type_id, '社交信息备注最大长度', '127', 1, 1, '社交信息备注最大127个中文字符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 8. 社交信息状态字典：blog_social_info_status
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('社交信息状态', 'blog_social_info_status', '社交信息启用/禁用状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_social_info_status_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_social_info_status' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_social_info_status_type_id, '启用', '1', 1, 1, '社交信息启用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_social_info_status_type_id, '禁用', '2', 2, 1, '社交信息禁用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 9. 文章置顶最大数量配置：blog_article_top_max_count
-- ============================================
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('文章置顶最大数量', 'blog_article_top_max_count', '博客文章最多可置顶的数量', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @blog_article_top_max_count_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'blog_article_top_max_count' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@blog_article_top_max_count_type_id, '文章置顶最大数量', '1', 1, 1, '默认最多置顶1篇文章', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
