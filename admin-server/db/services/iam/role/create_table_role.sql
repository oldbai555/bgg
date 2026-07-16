-- iam/role 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `admin_role` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
  `code` VARCHAR(64) NOT NULL COMMENT '角色编码（唯一）',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '角色描述',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1 启用，0 禁用',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_role_code` (`code`),
  KEY `idx_admin_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台角色表';
