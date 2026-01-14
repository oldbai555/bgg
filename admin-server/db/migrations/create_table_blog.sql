-- ============================================
-- 博客模块表结构 SQL
-- 模块：博客（标签、文章、审核）
-- 创建时间：2026-01-14
-- 说明：包含标签表、文章表、文章标签关联表、文章审核记录表
-- ============================================

-- 1. 文章标签表
CREATE TABLE IF NOT EXISTS `blog_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_tag_status）',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_tag_name` (`name`),
  KEY `idx_blog_tag_status` (`status`),
  KEY `idx_blog_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='博客标签表';

-- 2. 文章主表
CREATE TABLE IF NOT EXISTS `blog_article` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
  `content` LONGTEXT NOT NULL COMMENT '文章内容（Markdown 原文）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '文章状态（字典：blog_article_status）',
  `audit_status` TINYINT NOT NULL DEFAULT 1 COMMENT '审核状态（字典：blog_article_audit_status）',
  `cover` VARCHAR(512) DEFAULT '' COMMENT '封面图片 URL（为空时前端按标题首字生成占位封面）',
  `author_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '作者后台用户ID',
  `author_name` VARCHAR(64) DEFAULT '' COMMENT '作者姓名快照',
  `publish_time` BIGINT NOT NULL DEFAULT 0 COMMENT '上架时间(秒级时间戳,0表示未上架)',
  `summary` VARCHAR(512) DEFAULT '' COMMENT '文章摘要/简介（可选，公共列表优先使用）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_status` (`status`),
  KEY `idx_blog_article_audit_status` (`audit_status`),
  KEY `idx_blog_article_author_id` (`author_id`),
  KEY `idx_blog_article_publish_time` (`publish_time`),
  KEY `idx_blog_article_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='博客文章表';

-- 3. 文章-标签关联表（多对多）
CREATE TABLE IF NOT EXISTS `blog_article_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签 ID',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_tag_article_id` (`article_id`),
  KEY `idx_blog_article_tag_tag_id` (`tag_id`),
  KEY `idx_blog_article_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';

-- 4. 文章审核记录表
CREATE TABLE IF NOT EXISTS `blog_article_audit` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `audit_status` TINYINT NOT NULL COMMENT '本次审核结果（字典：blog_article_audit_status）',
  `audit_remark` VARCHAR(1000) DEFAULT '' COMMENT '审核意见',
  `auditor_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核人后台用户ID',
  `auditor_name` VARCHAR(64) DEFAULT '' COMMENT '审核人姓名快照',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_audit_article_id` (`article_id`),
  KEY `idx_blog_article_audit_auditor_id` (`auditor_id`),
  KEY `idx_blog_article_audit_created_at` (`created_at`),
  KEY `idx_blog_article_audit_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章审核记录表';

