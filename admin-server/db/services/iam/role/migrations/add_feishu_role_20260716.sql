-- iam/role 增量：新增"飞书"角色，供飞书扫码登录首次自动建号时默认分配
-- 幂等：角色按 code 唯一键 ON DUPLICATE KEY UPDATE；角色-权限关联按 NOT EXISTS 防重复插入

INSERT INTO `admin_role` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('飞书用户', 'feishu', '飞书扫码登录自动创建的用户默认分配此角色，权限需管理员按需追加', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

SET @feishu_role_id = (SELECT `id` FROM `admin_role` WHERE `code` = 'feishu' AND `deleted_at` = 0 LIMIT 1);

-- 默认权限：仅给"必备"的日常协作类只读权限，不含任何后台管理 CRUD/审计权限。
-- 仪表盘、在线聊天、消息通知管理等菜单本身未绑定权限（对所有登录用户默认放行），
-- 不需要在这里额外授权；只有绑定了权限编码的菜单才需要显式授予。
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
SELECT @feishu_role_id, p.`id`, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `admin_permission` p
WHERE p.`code` IN ('daily_short_sentence:list', 'file:list')
  AND NOT EXISTS (
    SELECT 1 FROM `admin_role_permission` rp
    WHERE rp.`role_id` = @feishu_role_id AND rp.`permission_id` = p.`id`
  );
