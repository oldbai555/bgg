-- iam/login_log 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 依赖 iam/menu 模块已初始化（/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system' AND `deleted_at` = 0 LIMIT 1);

-- 登录日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '登录日志',
    '/system/login-log',
    'system/LoginLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    31, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/login-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志详情按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @login_menu_id,
    '登录日志 详情按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_detail_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @login_menu_id
    AND `name` = '登录日志 详情按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @login_menu_id,
    '登录日志 导出按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    2, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @login_menu_id
    AND `name` = '登录日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('登录日志列表', 'login_log:list', '查看登录日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志详情', 'login_log:detail', '查看登录日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志导出', 'login_log:export', '导出登录日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @login_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:list' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_detail_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_export_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:export' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('登录日志列表', 'GET', '/api/v1/login-logs', '获取登录日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志详情', 'GET', '/api/v1/login-logs/detail', '获取登录日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志导出', 'GET', '/api/v1/login-logs/export', '导出登录日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @login_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@login_list_permission_id, @login_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_detail_permission_id, @login_detail_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_export_permission_id, @login_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 登录日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@login_list_permission_id, @login_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_detail_permission_id, @login_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_export_permission_id, @login_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联登录日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @login_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @login_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @login_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

