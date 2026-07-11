-- iam/audit_log 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `audit_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `audit_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '审计类型：permission_assign/role_change/config_modify/data_delete等',
  `audit_object` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '审计对象（模块/表名，如user_role/role_permission/role/config）',
  `audit_detail` TEXT COMMENT '审计详情（JSON格式，记录变更前后的数据）',
  `ip_address` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'IP 地址',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户代理',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_audit_log_user_id` (`user_id`),
  KEY `idx_audit_log_audit_type` (`audit_type`),
  KEY `idx_audit_log_audit_object` (`audit_object`),
  KEY `idx_audit_log_created_at` (`created_at`),
  KEY `idx_audit_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';
