-- iam/role 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 初始化角色：1=super_admin 超级管理员角色，2=admin 业务管理员角色
INSERT INTO `admin_role` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, '超级管理员', 'super_admin', '系统内置最高权限角色，拥有全部权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'admin', 'admin', '系统内置业务管理员角色，示例账号使用', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 飞书用户角色：飞书扫码登录自动创建的用户默认分配此角色，权限需管理员按需追加；
-- 默认权限（daily_short_sentence:list、file:list）在 iam/daily_short_sentence/
-- init_daily_short_sentence.sql 末尾赋予——那两个权限分别来自 iam/permission 和
-- iam/daily_short_sentence 模块，daily_short_sentence 是 init-dev-db.sh 里 iam 域最后
-- 一个模块，只有到那时两个权限才都已存在，本文件（role 模块，先于两者初始化）不能直接赋权。
INSERT INTO `admin_role` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('飞书用户', 'feishu', '飞书扫码登录自动创建的用户默认分配此角色，权限需管理员按需追加', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();
