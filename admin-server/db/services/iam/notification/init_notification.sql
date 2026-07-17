-- iam/notification 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 为初始化用户创建公告通知（公告已发布，需要给所有用户创建通知）
-- 注意：这里只给初始化时已存在的用户创建通知，后续新增的用户会在登录时自动获取未读公告
INSERT INTO `admin_notification` (`user_id`, `source_type`, `source_id`, `title`, `content`, `read_status`, `read_at`, `created_at`, `updated_at`, `deleted_at`)
SELECT
  1, -- 超级管理员（id=1）
  'notice',
  1, -- 公告ID
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  1, -- 未读（字典值：1=未读，2=已读）
  0,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_notification`
  WHERE `user_id` = 1 AND `source_type` = 'notice' AND `source_id` = 1 AND `deleted_at` = 0
);

INSERT INTO `admin_notification` (`user_id`, `source_type`, `source_id`, `title`, `content`, `read_status`, `read_at`, `created_at`, `updated_at`, `deleted_at`)
SELECT
  2, -- admin业务管理员（id=2）
  'notice',
  1, -- 公告ID
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  1, -- 未读（字典值：1=未读，2=已读）
  0,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_notification`
  WHERE `user_id` = 2 AND `source_type` = 'notice' AND `source_id` = 1 AND `deleted_at` = 0
);

-- 依赖 iam/menu 模块已初始化（/admin/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system' AND `deleted_at` = 0 LIMIT 1);

-- 消息通知管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@system_menu_id, '消息通知管理', '/admin/system/notification', 'system/NotificationList', 'ele-Bell', 2, 22, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 消息通知管理接口（只需要AuthMiddleware，不需要PermissionMiddleware）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('消息通知管理列表', 'GET', '/api/v1/notifications', '获取消息通知管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知全部已读', 'PUT', '/api/v1/notifications/read-all', '标记所有消息通知为已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知标记已读', 'PUT', '/api/v1/notifications/read', '标记单个消息通知为已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知清除已读', 'DELETE', '/api/v1/notifications/read', '清除所有已读消息通知', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知删除', 'DELETE', '/api/v1/notifications', '删除消息通知', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

