-- iam/login_log 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_login_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `ip_address` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登录 IP 地址',
  `location` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '登录地点（通过IP解析）',
  `browser` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '浏览器',
  `os` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '操作系统',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户代理',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '登录状态（字典 login_status）：1 成功，2 失败',
  `message` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '登录消息（失败原因或成功提示）',
  `login_at` BIGINT NOT NULL DEFAULT 0 COMMENT '登录时间(秒级时间戳)',
  `logout_at` BIGINT NOT NULL DEFAULT 0 COMMENT '登出时间(秒级时间戳,0表示未登出)',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_login_log_user_id` (`user_id`),
  KEY `idx_login_log_username` (`username`),
  KEY `idx_login_log_status` (`status`),
  KEY `idx_login_log_login_at` (`login_at`),
  KEY `idx_login_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';
