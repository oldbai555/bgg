-- iam/config 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_config` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `group` VARCHAR(64) NOT NULL DEFAULT 'default' COMMENT '配置分组（如 system、app、theme 等）',
  `key` VARCHAR(128) NOT NULL COMMENT '配置键（唯一，格式：group:key）',
  `value` TEXT COMMENT '配置值（JSON 格式存储复杂数据）',
  `type` VARCHAR(32) NOT NULL DEFAULT 'string' COMMENT '配置类型（string、number、boolean、json）',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置描述',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_config_key` (`key`),
  KEY `idx_admin_config_group` (`group`),
  KEY `idx_admin_config_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
