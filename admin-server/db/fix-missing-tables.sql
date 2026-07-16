-- iam/user_third_party 建表
-- 第三方登录账号绑定表（飞书扫码登录等），不修改 admin_user 表结构，
-- 通过 user_id 关联回 admin_user.id，为后续接入更多第三方登录方式（企业微信/钉钉等）留出扩展空间

CREATE TABLE IF NOT EXISTS `admin_user_third_party` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '关联 admin_user.id',
  `provider` VARCHAR(32) NOT NULL COMMENT '第三方登录渠道标识，如 feishu',
  `open_id` VARCHAR(128) NOT NULL COMMENT '第三方渠道下的用户唯一标识',
  `union_id` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '第三方渠道下的开发者主体维度唯一标识',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_user_third_party_provider_open_id` (`provider`, `open_id`),
  KEY `idx_admin_user_third_party_user_id` (`user_id`),
  KEY `idx_admin_user_third_party_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方登录账号绑定表';
CREATE TABLE IF NOT EXISTS `demo` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(255) NOT NULL COMMENT '描述',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1 启用，0 禁用',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_demo_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='演示功能表';

