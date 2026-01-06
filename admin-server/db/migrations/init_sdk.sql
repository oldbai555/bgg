-- sdk模块初始化 SQL
-- 功能组: sdk
-- 功能名称: SDK 管理（API Key / 调用记录 / 接口管理）

-- ============================================
-- 1. 获取父目录 ID（默认放在根或指定父目录）
-- ============================================
SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk' AND `type` = 1 AND `deleted_at` = 0 LIMIT 1),
  0
);

-- ============================================
-- 2. 插入菜单数据（目录 + 三个子菜单）
-- ============================================
-- SDK 管理主目录
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    'SDK管理',
    '/sdk',
    '',
    'ele-Key',
    1, -- 目录
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_root_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk' AND `type` = 1 AND `deleted_at` = 0 LIMIT 1);

-- API Key 管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @sdk_root_id,
    'API管理',
    '/sdk/api-key',
    'sdk/ApiKeyList',
    'ele-List',
    2, -- 菜单
    1,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_api_key_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/api-key' AND `deleted_at` = 0 LIMIT 1);

-- 调用记录菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @sdk_root_id,
    '调用记录',
    '/sdk/call-log',
    'sdk/SdkCallLogList',
    'ele-DocumentChecked',
    2,
    2,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_call_log_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/call-log' AND `deleted_at` = 0 LIMIT 1);

-- SDK 接口管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @sdk_root_id,
    'SDK接口管理',
    '/sdk/interface',
    'sdk/SdkInterfaceList',
    'ele-Setting',
    2,
    3,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_interface_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/interface' AND `deleted_at` = 0 LIMIT 1);

-- 按钮/操作（API Key）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_api_key_menu_id, 'API Key 新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 分配接口', '', '', '', 3, 4, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @btn_api_key_create = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 新增' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_update = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 编辑' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_delete = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 删除' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_bind   = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 分配接口' AND `deleted_at`=0 LIMIT 1);

-- 按钮/操作（调用记录）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_call_log_menu_id, '调用记录导出', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @btn_call_log_export = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_call_log_menu_id AND `name`='调用记录导出' AND `deleted_at`=0 LIMIT 1);

-- 按钮/操作（接口管理）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_interface_menu_id, 'SDK接口新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_interface_menu_id, 'SDK接口编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_interface_menu_id, 'SDK接口删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @btn_interface_create = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口新增' AND `deleted_at`=0 LIMIT 1);
SET @btn_interface_update = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口编辑' AND `deleted_at`=0 LIMIT 1);
SET @btn_interface_delete = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口删除' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 3. 插入权限数据
-- ============================================
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK API Key 列表', 'sdk:key:list', '查看 SDK API Key 列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 新增', 'sdk:key:create', '新增 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 编辑', 'sdk:key:update', '编辑 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 删除', 'sdk:key:delete', '删除 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 分配接口', 'sdk:key:bind_api', '分配接口给 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录列表', 'sdk:call_log:list', '查看 SDK 调用记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录导出', 'sdk:call_log:export', '导出 SDK 调用记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口列表', 'sdk:interface:list', '查看 SDK 接口列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口新增', 'sdk:interface:create', '新增 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口编辑', 'sdk:interface:update', '编辑 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口删除', 'sdk:interface:delete', '删除 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @perm_key_list      = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_create    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:create' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_update    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:update' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_delete    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:delete' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_bind      = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:bind_api' AND `deleted_at`=0 LIMIT 1);
SET @perm_call_list     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:call_log:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_call_export   = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:call_log:export' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_list       = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_create     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:create' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_update     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:update' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_delete     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:delete' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 4. 插入接口数据（仅管理端）
-- ============================================
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK API Key 列表', 'GET', '/api/v1/sdk/key/list', 'SDK API Key 列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 新增', 'POST', '/api/v1/sdk/key/create', '新增 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 编辑', 'POST', '/api/v1/sdk/key/update', '编辑 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 删除', 'POST', '/api/v1/sdk/key/delete', '删除 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 接口列表', 'GET', '/api/v1/sdk/key/apis', '查看指定 key 的接口授权', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 接口保存', 'POST', '/api/v1/sdk/key/apis/save', '为 key 保存接口授权与限频', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录列表', 'GET', '/api/v1/sdk/call/log/list', 'SDK 调用记录列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录导出', 'GET', '/api/v1/sdk/call/log/export', '导出 SDK 调用记录', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口列表', 'GET', '/api/v1/sdk/interface/list', 'SDK 接口列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口新增', 'POST', '/api/v1/sdk/interface/create', '新增 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口编辑', 'POST', '/api/v1/sdk/interface/update', '编辑 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口删除', 'POST', '/api/v1/sdk/interface/delete', '删除 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @api_key_list     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/list' AND `deleted_at`=0 LIMIT 1);
SET @api_key_create   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/create' AND `deleted_at`=0 LIMIT 1);
SET @api_key_update   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/update' AND `deleted_at`=0 LIMIT 1);
SET @api_key_delete   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/delete' AND `deleted_at`=0 LIMIT 1);
SET @api_key_apis     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/apis' AND `deleted_at`=0 LIMIT 1);
SET @api_key_save     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/apis/save' AND `deleted_at`=0 LIMIT 1);
SET @api_call_list    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/call/log/list' AND `deleted_at`=0 LIMIT 1);
SET @api_call_export  = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/call/log/export' AND `deleted_at`=0 LIMIT 1);
SET @api_if_list      = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/list' AND `deleted_at`=0 LIMIT 1);
SET @api_if_create    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/create' AND `deleted_at`=0 LIMIT 1);
SET @api_if_update    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/update' AND `deleted_at`=0 LIMIT 1);
SET @api_if_delete    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/delete' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 5. 权限-菜单关联
-- ============================================
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@perm_key_list, @sdk_api_key_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_create, @btn_api_key_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_update, @btn_api_key_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_delete, @btn_api_key_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @btn_api_key_bind, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_list, @sdk_call_log_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_export, @btn_call_log_export, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_list, @sdk_interface_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_create, @btn_interface_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_update, @btn_interface_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_delete, @btn_interface_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 6. 权限-接口关联
-- ============================================
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@perm_key_list, @api_key_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_create, @api_key_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_update, @api_key_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_delete, @api_key_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @api_key_apis, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @api_key_save, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_list, @api_call_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_export, @api_call_export, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_list, @api_if_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_create, @api_if_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_update, @api_if_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_delete, @api_if_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

