-- iam/role 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 初始化角色：1=super_admin 超级管理员角色，2=admin 业务管理员角色
INSERT INTO `admin_role` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, '超级管理员', 'super_admin', '系统内置最高权限角色，拥有全部权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'admin', 'admin', '系统内置业务管理员角色，示例账号使用', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;
