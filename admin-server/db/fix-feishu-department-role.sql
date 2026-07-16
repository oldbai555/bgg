-- 补齐飞书登录自动建号的默认部门/角色分配（见 commit 4543415f 的说明），oldbai 库之前
-- 没跑过这两个迁移。全部幂等，重复执行安全。

-- 1) 飞书待分配部门（db/services/iam/department/migrations/add_feishu_pending_department_20260716.sql）
INSERT INTO `admin_department` (`id`, `parent_id`, `name`, `order_num`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (2, 0, '飞书待分配', 99, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at` = 0, `name` = '飞书待分配';

-- 2) 飞书角色 + 默认权限（db/services/iam/role/migrations/add_feishu_role_20260716.sql）
INSERT INTO `admin_role` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('飞书用户', 'feishu', '飞书扫码登录自动创建的用户默认分配此角色，权限需管理员按需追加', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

SET @feishu_role_id = (SELECT `id` FROM `admin_role` WHERE `code` = 'feishu' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
SELECT @feishu_role_id, p.`id`, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `admin_permission` p
WHERE p.`code` IN ('daily_short_sentence:list', 'file:list')
  AND NOT EXISTS (
    SELECT 1 FROM `admin_role_permission` rp
    WHERE rp.`role_id` = @feishu_role_id AND rp.`permission_id` = p.`id`
  );

-- 3) 存量补齐：oldbai 库里已经因为这两个迁移没跑过、建号时漏分配的飞书用户
--    （本次上线联调期间产生的账号，见 admin_user.username LIKE 'feishu_%'）
UPDATE `admin_user`
SET `department_id` = 2
WHERE `username` LIKE 'feishu\_%' AND `department_id` = 0;

INSERT INTO `admin_user_role` (`user_id`, `role_id`, `created_at`, `updated_at`)
SELECT u.`id`, @feishu_role_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `admin_user` u
WHERE u.`username` LIKE 'feishu\_%'
  AND NOT EXISTS (
    SELECT 1 FROM `admin_user_role` ur WHERE ur.`user_id` = u.`id` AND ur.`role_id` = @feishu_role_id
  );
