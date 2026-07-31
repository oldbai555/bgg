-- iam/audit_log 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 依赖 iam/menu 模块已初始化（/admin/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system' AND `deleted_at` = 0 LIMIT 1);

-- 审计日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '审计日志',
    '/admin/system/audit-log',
    'monitoring/AuditLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    32, -- 排序值
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
SET @audit_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/admin/system/audit-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @audit_menu_id,
    '审计日志 导出按钮',
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
SET @audit_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @audit_menu_id
    AND `name` = '审计日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('审计日志列表', 'audit_log:list', '查看审计日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志详情', 'audit_log:detail', '查看审计日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志导出', 'audit_log:export', '导出审计日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @audit_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:list' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_detail_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_export_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:export' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('审计日志列表', 'GET', '/api/v1/audit-logs', '获取审计日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志详情', 'GET', '/api/v1/audit-logs/detail', '获取审计日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志导出', 'GET', '/api/v1/audit-logs/export', '导出审计日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @audit_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@audit_list_permission_id, @audit_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_export_permission_id, @audit_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 审计日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@audit_list_permission_id, @audit_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_detail_permission_id, @audit_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_export_permission_id, @audit_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联审计日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @audit_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @audit_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @audit_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

