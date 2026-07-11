-- iam/performance_log 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_performance_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `method` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '请求方法：GET/POST/PUT/DELETE等',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径',
  `status_code` INT NOT NULL DEFAULT 0 COMMENT '响应状态码',
  `duration` INT NOT NULL DEFAULT 0 COMMENT '请求耗时（毫秒）',
  `is_slow` TINYINT NOT NULL DEFAULT 0 COMMENT '是否慢接口：0 否，1 是',
  `slow_threshold` INT NOT NULL DEFAULT 0 COMMENT '慢接口阈值（毫秒）',
  `ip_address` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'IP 地址',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户代理',
  `error_msg` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '错误信息（状态码>=400时记录）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_performance_log_path` (`path`),
  KEY `idx_performance_log_created_at` (`created_at`),
  KEY `idx_performance_log_is_slow` (`is_slow`),
  KEY `idx_performance_log_duration` (`duration`),
  KEY `idx_performance_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='接口性能监控日志表';
