-- iam/notification 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_notification` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `source_type` VARCHAR(32) NOT NULL COMMENT '消息来源类型（通过字典定义：chat、notice等）',
  `source_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源ID（如公告ID、聊天消息ID等）',
  `title` VARCHAR(255) NOT NULL COMMENT '消息标题',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `read_status` TINYINT NOT NULL DEFAULT 1 COMMENT '已读状态（字典 read_status）：1 未读，2 已读',
  `read_at` BIGINT NOT NULL DEFAULT 0 COMMENT '已读时间(秒级时间戳)',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_admin_notification_user_id` (`user_id`),
  KEY `idx_admin_notification_source_type` (`source_type`),
  KEY `idx_admin_notification_read_status` (`read_status`),
  KEY `idx_admin_notification_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息通知管理表';
