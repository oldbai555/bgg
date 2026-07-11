-- ============================================
-- 博客扩展功能表结构 SQL
-- 模块：博客扩展功能（友情链接、社交信息、文章置顶）
-- 创建时间：2026-01-16
-- 说明：包含友情链接表、社交信息表，以及文章表置顶字段扩展
-- ============================================

-- ============================================
-- 1. 友情链接表
-- ============================================
CREATE TABLE IF NOT EXISTS `blog_friend_link` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '链接名称（最大15个中文字符，通过字典配置）',
  `url` VARCHAR(255) NOT NULL COMMENT '目标链接（最大255字符，通过字典配置）',
  `remark` VARCHAR(512) DEFAULT '' COMMENT '备注（最大127个中文字符，通过字典配置）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_friend_link_status）',
  `order_num` INT NOT NULL DEFAULT 0 COMMENT '排序值（数字越小越靠前）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_friend_link_status` (`status`),
  KEY `idx_blog_friend_link_order_num` (`order_num`),
  KEY `idx_blog_friend_link_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='友情链接表';

-- ============================================
-- 2. 社交信息表
-- ============================================
CREATE TABLE IF NOT EXISTS `blog_social_info` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '社交平台名称（最大15个中文字符，通过字典配置）',
  `url` VARCHAR(255) NOT NULL COMMENT '目标链接（最大255字符，通过字典配置）',
  `remark` VARCHAR(512) DEFAULT '' COMMENT '备注（最大127个中文字符，通过字典配置）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_social_info_status）',
  `order_num` INT NOT NULL DEFAULT 0 COMMENT '排序值（数字越小越靠前）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_social_info_status` (`status`),
  KEY `idx_blog_social_info_order_num` (`order_num`),
  KEY `idx_blog_social_info_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社交信息表';

-- ============================================
-- 3. 文章表扩展（添加置顶字段）
-- ============================================
-- 检查字段是否已存在，如果不存在则添加
SET @dbname = DATABASE();
SET @tablename = 'blog_article';
SET @columnname = 'is_top';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname)
  ) > 0,
  'SELECT 1',
  CONCAT('ALTER TABLE `', @tablename, '` ADD COLUMN `', @columnname, '` TINYINT NOT NULL DEFAULT 0 COMMENT ''是否置顶：0=否，1=是'' AFTER `summary`, ADD KEY `idx_blog_article_is_top` (`is_top`)')
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;
