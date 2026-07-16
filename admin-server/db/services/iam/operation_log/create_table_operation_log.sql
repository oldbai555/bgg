-- iam/operation_log 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_operation_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `operation_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '操作类型：create/update/delete/query/export等',
  `operation_object` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '操作对象（模块/表名，如user/role/permission）',
  `method` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '请求方法：GET/POST/PUT/DELETE',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径',
  `request_params` TEXT COMMENT '请求参数（JSON格式）',
  `response_code` INT NOT NULL DEFAULT 0 COMMENT '响应状态码',
  `response_msg` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '响应消息',
  `ip_address` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'IP 地址',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户代理',
  `duration` INT NOT NULL DEFAULT 0 COMMENT '请求耗时（毫秒）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_operation_log_user_id` (`user_id`),
  KEY `idx_operation_log_operation_type` (`operation_type`),
  KEY `idx_operation_log_operation_object` (`operation_object`),
  KEY `idx_operation_log_created_at` (`created_at`),
  KEY `idx_operation_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
