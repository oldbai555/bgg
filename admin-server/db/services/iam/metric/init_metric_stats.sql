-- ============================================
-- 数据统计模块初始化 SQL
-- 功能组: metric_admin
-- 功能名称: 数据统计（PV/UV/VV/IP）
-- 创建时间: 2025-01-15
-- ============================================

-- ============================================
-- 1. 获取系统管理目录 ID（parent_id=2）
-- ============================================
SET @system_dir_id = (SELECT `id` FROM `admin_menu` WHERE `id` = 2 AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 2. 插入菜单数据
-- ============================================
-- 数据统计主菜单（放在系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_dir_id,
    '数据统计',
    '/system/metric-stats',
    'system/MetricStats',
    'ele-DataAnalysis',
    2, -- 类型：2 菜单
    21, -- 排序值（放在系统管理下，在任务管理之后）
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @main_menu_id = LAST_INSERT_ID();

-- ============================================
-- 3. 插入权限数据
-- ============================================
-- 数据统计查询权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '数据统计',
    'metric:stats',
    '查看PV/UV/VV/IP统计数据',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @stats_permission_id = LAST_INSERT_ID();

-- ============================================
-- 4. 插入接口数据（注意：接口可能已通过路由同步自动创建，这里使用ON DUPLICATE KEY UPDATE）
-- ============================================
-- 数据统计接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '数据统计',
    'GET',
    '/api/v1/metrics/stats',
    '获取PV/UV/VV/IP统计数据',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `description`=VALUES(`description`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;

SET @stats_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/metrics/stats' AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 5. 插入权限-菜单关联数据
-- ============================================
-- 数据统计权限 -> 数据统计主菜单
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@stats_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 6. 插入权限-接口关联数据
-- ============================================
-- 数据统计权限 -> GET /api/v1/metrics/stats接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@stats_permission_id, @stats_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();
