-- iam/operation_log 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 依赖 iam/menu 模块已初始化（/admin/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system' AND `deleted_at` = 0 LIMIT 1);

-- 操作日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '操作日志',
    '/admin/system/operation-log',
    'monitoring/OperationLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    30, -- 排序值
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
SET @operation_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/admin/system/operation-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @operation_menu_id,
    '操作日志 导出按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否（按钮不显示在菜单中）
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
SET @operation_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @operation_menu_id
    AND `name` = '操作日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('操作日志列表', 'operation_log:list', '查看操作日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志详情', 'operation_log:detail', '查看操作日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志导出', 'operation_log:export', '导出操作日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @operation_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'operation_log:list' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_detail_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'operation_log:detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_export_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'operation_log:export' AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('操作日志列表', 'GET', '/api/v1/operation-logs', '获取操作日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志详情', 'GET', '/api/v1/operation-logs/detail', '获取操作日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志导出', 'GET', '/api/v1/operation-logs/export', '导出操作日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @operation_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@operation_list_permission_id, @operation_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_export_permission_id, @operation_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 操作日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@operation_list_permission_id, @operation_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_detail_permission_id, @operation_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_export_permission_id, @operation_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联操作日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @operation_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @operation_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @operation_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

