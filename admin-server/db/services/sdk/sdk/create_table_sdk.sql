-- sdk/sdk 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `sdk_key` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',
  `api_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '对外 API Key（唯一）',
  `api_secret` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '对外 API Secret（唯一）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `expire_at` BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间(秒级时间戳)，0 表示长期有效',
  `ip_whitelist` TEXT NOT NULL COMMENT 'IP 白名单，逗号或换行分隔，空表示不限制',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_key_api_key` (`api_key`),
  UNIQUE KEY `uk_sdk_key_api_secret` (`api_secret`),
  KEY `idx_sdk_key_status` (`status`),
  KEY `idx_sdk_key_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK Key';

CREATE TABLE IF NOT EXISTS `sdk_interface` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口名称',
  `api_code` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口编码（唯一，便于绑定与校验）',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径（对外 SDK 路由）',
  `method` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'HTTP 方法',
  `rate_limit_default` INT NOT NULL DEFAULT 0 COMMENT '默认限频上限（0 表示不限，实际逻辑读取字典兜底）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_interface_code` (`api_code`),
  KEY `idx_sdk_interface_status` (`status`),
  KEY `idx_sdk_interface_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK 对外接口';

CREATE TABLE IF NOT EXISTS `sdk_key_api` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `sdk_key_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK Key ID',
  `sdk_interface_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK 接口 ID',
  `custom_rate_limit` INT NOT NULL DEFAULT 0 COMMENT '自定义限频上限（0 表示使用接口默认值或字典默认值）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_key_api` (`sdk_key_id`, `sdk_interface_id`),
  KEY `idx_sdk_key_api_key_id` (`sdk_key_id`),
  KEY `idx_sdk_key_api_interface_id` (`sdk_interface_id`),
  KEY `idx_sdk_key_api_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK Key 接口授权';

CREATE TABLE IF NOT EXISTS `sdk_call_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `sdk_key_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK Key ID',
  `sdk_interface_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK 接口 ID',
  `api_code` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口编码快照',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径',
  `method` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'HTTP 方法',
  `ip` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '调用 IP',
  `user_agent` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'User-Agent',
  `req_body` MEDIUMTEXT COMMENT '请求体（脱敏/截断后写入）',
  `resp_body` MEDIUMTEXT COMMENT '响应体（脱敏/截断后写入）',
  `resp_code` INT NOT NULL DEFAULT 0 COMMENT 'HTTP 状态码',
  `duration_ms` INT NOT NULL DEFAULT 0 COMMENT '耗时（毫秒）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_sdk_call_log_key` (`sdk_key_id`),
  KEY `idx_sdk_call_log_interface` (`sdk_interface_id`),
  KEY `idx_sdk_call_log_created_at` (`created_at`),
  KEY `idx_sdk_call_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK 调用日志';
