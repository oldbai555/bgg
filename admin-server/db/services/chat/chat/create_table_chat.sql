-- chat/chat 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `chat` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '聊天名称（群组名称，私聊为空）',
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '聊天类型：1私聊，2群组',
  `avatar` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像URL（群组头像，私聊为空）',
  `description` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '描述（群组描述，私聊为空）',
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID（群组创建人，私聊为0）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_chat_type` (`type`),
  KEY `idx_chat_created_by` (`created_by`),
  KEY `idx_chat_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天表（私聊、群组）';

CREATE TABLE IF NOT EXISTS `chat_user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `chat_id` BIGINT UNSIGNED NOT NULL COMMENT '聊天 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
  `joined_at` BIGINT NOT NULL DEFAULT 0 COMMENT '加入时间(秒级时间戳)',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chat_user` (`chat_id`,`user_id`),
  KEY `idx_chat_user_chat_id` (`chat_id`),
  KEY `idx_chat_user_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天-用户关联表';

CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `chat_id` BIGINT UNSIGNED NOT NULL COMMENT '聊天 ID（关联chat表）',
  `from_user_id` BIGINT UNSIGNED NOT NULL COMMENT '发送用户 ID',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `message_type` TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型：1文本，2图片，3文件',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1已发送，2已读，3已撤回',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_chat_message_chat_id` (`chat_id`),
  KEY `idx_chat_message_from_user_id` (`from_user_id`),
  KEY `idx_chat_message_created_at` (`created_at`),
  KEY `idx_chat_message_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';
