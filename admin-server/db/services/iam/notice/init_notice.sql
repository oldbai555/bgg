-- iam/notice 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 插入一条欢迎公告（已发布状态）
INSERT INTO `admin_notice` (`id`, `title`, `content`, `type`, `status`, `publish_time`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  1,
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  1, -- 类型：1 普通公告
  2, -- 状态：2 已发布
  UNIX_TIMESTAMP(), -- 发布时间：当前时间
  1, -- 创建人：超级管理员（id=1）
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `content`=VALUES(`content`),
  `type`=VALUES(`type`),
  `status`=VALUES(`status`),
  `publish_time`=VALUES(`publish_time`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;


-- 依赖 iam/menu 模块已初始化（/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@system_menu_id, '公告管理', '/system/notice', 'system/NoticeList', 'ele-Document', 2, 21, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system/notice' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@notice_menu_id, '公告管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_menu_id, '公告管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_menu_id, '公告管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 新增按钮' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 编辑按钮' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('公告管理列表', 'notice:list', '查看公告管理列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理新增', 'notice:create', '新增公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理编辑', 'notice:update', '编辑公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理删除', 'notice:delete', '删除公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:list' AND `deleted_at` = 0 LIMIT 1);
SET @notice_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:create' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:update' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:delete' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('公告管理列表', 'GET', '/api/v1/notices', '获取公告管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理新增', 'POST', '/api/v1/notices', '新增公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理编辑', 'PUT', '/api/v1/notices', '编辑公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理删除', 'DELETE', '/api/v1/notices', '删除公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_list_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_create_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'PUT' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@notice_list_permission_id, @notice_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_create_permission_id, @notice_create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_update_permission_id, @notice_update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_delete_permission_id, @notice_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 公告管理 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@notice_list_permission_id, @notice_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_create_permission_id, @notice_create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_update_permission_id, @notice_update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_delete_permission_id, @notice_delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

