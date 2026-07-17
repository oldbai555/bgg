-- chat/chat 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 创建默认企业群组（id=1，type=2群组，created_by=1超级管理员）
INSERT INTO `chat` (`id`, `name`, `type`, `avatar`, `description`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  1,
  '企业群组',
  2, -- 类型：2 群组
  '',
  '默认企业群组，所有用户自动加入',
  1, -- 创建人：超级管理员（id=1）
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `type`=VALUES(`type`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 将初始的两个用户（id=1和id=2）都加入默认企业群组
INSERT INTO `chat_user` (`id`, `chat_id`, `user_id`, `joined_at`, `created_at`, `updated_at`)
VALUES
  (1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- 超级管理员加入群组
  (2, 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- 业务管理员加入群组
ON DUPLICATE KEY UPDATE
  `joined_at`=VALUES(`joined_at`),
  `updated_at`=UNIX_TIMESTAMP();

-- 创建用户1和用户2之间的私聊（id=2，type=1私聊）
INSERT INTO `chat` (`id`, `name`, `type`, `avatar`, `description`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  2,
  '', -- 私聊名称为空，前端根据对方用户信息显示
  1, -- 类型：1 私聊
  '',
  '',
  0, -- 私聊创建人为0
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE
  `type`=VALUES(`type`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 将用户1和用户2加入私聊
INSERT INTO `chat_user` (`id`, `chat_id`, `user_id`, `joined_at`, `created_at`, `updated_at`)
VALUES
  (3, 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- 用户1加入私聊
  (4, 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- 用户2加入私聊
ON DUPLICATE KEY UPDATE
  `joined_at`=VALUES(`joined_at`),
  `updated_at`=UNIX_TIMESTAMP();

-- 聊天室目录
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (0, '聊天室', '/admin/chatroom', '', 'ele-ChatDotRound', 1, 20, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chatroom_dir_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/chatroom' AND `deleted_at` = 0 LIMIT 1);

-- 在线聊天菜单（无需权限，只要登录就可以访问）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '在线聊天', '/admin/chatroom/chat', 'chat/ChatList', 'ele-ChatLineRound', 2, 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 聊天记录管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '聊天记录管理', '/admin/chatroom/chat-message', 'chat/ChatMessageList', 'ele-Document', 2, 2, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/chatroom/chat-message' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理删除按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chat_message_menu_id, '聊天记录管理 删除按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_message_menu_id AND `name` = '聊天记录管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('聊天记录列表', 'chat_message:list', '查看聊天记录列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('聊天记录删除', 'chat_message:delete', '删除聊天记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat_message:list' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat_message:delete' AND `deleted_at` = 0 LIMIT 1);

-- 在线聊天接口（无需权限，只需要AuthMiddleware）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('聊天消息发送', 'POST', '/api/v1/chats/messages', '发送聊天消息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('可聊天用户列表', 'GET', '/api/v1/chats/users', '获取可聊天用户列表（包含部门-角色-昵称信息）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 聊天记录管理接口（需要权限）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('聊天记录列表', 'GET', '/api/v1/chats/messages', '获取聊天记录列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('聊天记录删除', 'DELETE', '/api/v1/chats/messages', '删除聊天记录', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_list_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/chats/messages' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/chats/messages' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@chat_message_list_permission_id, @chat_message_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_message_delete_permission_id, @chat_message_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 群组管理菜单（在聊天室目录下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '群组管理', '/admin/chatroom/chat-group', 'chat/ChatGroupList', 'ele-ChatDotRound', 2, 3, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_group_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/chatroom/chat-group' AND `deleted_at` = 0 LIMIT 1);

-- 群组管理按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@chat_group_menu_id, '群组管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_group_menu_id, '群组管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_group_menu_id, '群组管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_group_create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 新增按钮' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 编辑按钮' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 群组管理 权限-菜单 关联
SET @chat_group_detail_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:detail' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:create' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:update' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:delete' AND `deleted_at` = 0 LIMIT 1);
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@chat_group_detail_permission_id, @chat_group_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_create_permission_id, @chat_group_create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_update_permission_id, @chat_group_update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_delete_permission_id, @chat_group_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 聊天记录管理 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@chat_message_list_permission_id, @chat_message_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_message_delete_permission_id, @chat_message_delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

