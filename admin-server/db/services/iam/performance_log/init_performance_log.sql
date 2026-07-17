-- iam/performance_log 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 依赖 iam/menu 模块已初始化（/admin/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system' AND `deleted_at` = 0 LIMIT 1);

-- 性能监控日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '性能监控日志',
    '/admin/system/performance-log',
    'monitoring/PerformanceLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    33, -- 排序值
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
SET @performance_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/admin/system/performance-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志列表权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '性能监控日志列表',
    'performance_log:list',
    '查看性能监控日志列表',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @performance_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'performance_log:list' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志列表接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '性能监控日志列表',
    'GET',
    '/api/v1/performance-logs',
    '获取性能监控日志列表',
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @performance_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/performance-logs' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@performance_list_permission_id, @performance_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 性能监控日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@performance_list_permission_id, @performance_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联性能监控日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES (1, @performance_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

