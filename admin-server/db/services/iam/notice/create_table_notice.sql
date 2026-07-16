-- iam/notice 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_notice` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `title` VARCHAR(255) NOT NULL COMMENT '公告标题',
  `content` TEXT NOT NULL COMMENT '公告内容',
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '公告类型：1 普通公告，2 重要公告，3 紧急公告',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1 草稿，2 已发布',
  `publish_time` BIGINT NOT NULL DEFAULT 0 COMMENT '发布时间(秒级时间戳)',
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_admin_notice_type` (`type`),
  KEY `idx_admin_notice_status` (`status`),
  KEY `idx_admin_notice_publish_time` (`publish_time`),
  KEY `idx_admin_notice_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告管理表';
