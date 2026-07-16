-- iam/monitor 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 依赖 iam/menu 模块已初始化（/system 目录菜单已存在）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system' AND `deleted_at` = 0 LIMIT 1);

-- 系统监控主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '系统监控',
    '/system/monitor',
    'system/MonitorList',
    'ele-Monitor',
    2, -- 类型：2 菜单
    34, -- 排序值
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
SET @monitor_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/monitor' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控查看权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '系统监控查看',
    'monitor:view',
    '查看系统监控信息',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @monitor_view_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'monitor:view' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('系统监控状态', 'GET', '/api/v1/monitor/status', '获取系统资源使用情况（CPU、内存、磁盘、网络）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('系统统计', 'GET', '/api/v1/monitor/stats', '获取系统统计数据（用户数、角色数、权限数等）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @monitor_status_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/monitor/status' AND `deleted_at` = 0
  LIMIT 1
);
SET @monitor_stats_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/monitor/stats' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@monitor_view_permission_id, @monitor_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 系统监控 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@monitor_view_permission_id, @monitor_status_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@monitor_view_permission_id, @monitor_stats_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联系统监控权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES (1, @monitor_view_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

