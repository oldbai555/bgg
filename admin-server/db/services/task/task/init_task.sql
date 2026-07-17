-- ============================================
-- 异步任务模块初始化 SQL
-- 功能组: task
-- 功能名称: 异步任务管理
-- 创建时间: 2025-01-15
-- ============================================

-- ============================================
-- 1. 获取系统管理目录 ID（parent_id=2）
-- ============================================
SET @system_dir_id = (SELECT `id` FROM `admin_menu` WHERE `id` = 2 AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 2. 插入菜单数据
-- ============================================
-- 异步任务管理主菜单（放在系统管理下）；admin_menu 没有唯一键，ON DUPLICATE KEY UPDATE
-- 对它不会触发，改用 INSERT ... SELECT ... WHERE NOT EXISTS 保证幂等
-- （同 content/blog/init_blog.sql）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @system_dir_id, '任务管理', '/admin/system/task', 'task/TaskList', 'ele-Document', 2, 20, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `path` = '/admin/system/task' AND `deleted_at` = 0);

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system/task' AND `deleted_at` = 0 LIMIT 1);

-- 任务管理查看按钮（用于查看任务详情）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @main_menu_id, '任务管理 查看按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @main_menu_id AND `name` = '任务管理 查看按钮' AND `deleted_at` = 0);

SET @view_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @main_menu_id AND `name` = '任务管理 查看按钮' AND `deleted_at` = 0 LIMIT 1);

-- 任务管理取消按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @main_menu_id, '任务管理 取消按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @main_menu_id AND `name` = '任务管理 取消按钮' AND `deleted_at` = 0);

SET @cancel_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @main_menu_id AND `name` = '任务管理 取消按钮' AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 3. 插入权限数据
-- ============================================
-- 任务管理列表权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('任务列表', 'task:list', '查看任务列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='task:list' AND `deleted_at`=0 LIMIT 1);

-- 任务管理详情权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('任务详情', 'task:detail', '查看任务详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @detail_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='task:detail' AND `deleted_at`=0 LIMIT 1);

-- 任务管理取消权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('任务取消', 'task:cancel', '取消任务', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @cancel_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='task:cancel' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 4. 插入接口数据（注意：接口可能已通过路由同步自动创建，这里使用ON DUPLICATE KEY UPDATE）
-- ============================================
-- 任务列表接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '任务列表',
    'GET',
    '/api/v1/tasks',
    '获取任务列表',
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

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/tasks' AND `deleted_at` = 0 LIMIT 1);

-- 任务详情接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '任务详情',
    'GET',
    '/api/v1/tasks/detail',
    '获取任务详情',
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `description`=VALUES(`description`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;

SET @detail_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/tasks/detail' AND `deleted_at` = 0 LIMIT 1);

-- 任务取消接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '任务取消',
    'POST',
    '/api/v1/tasks/cancel',
    '取消任务',
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `description`=VALUES(`description`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;

SET @cancel_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/tasks/cancel' AND `deleted_at` = 0 LIMIT 1);

-- 最近任务接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '最近任务',
    'GET',
    '/api/v1/tasks/recent',
    '获取最近任务列表（用于浮动球）',
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
    `name`=VALUES(`name`),
    `description`=VALUES(`description`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;

SET @recent_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/tasks/recent' AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 5. 插入权限-菜单关联数据
-- ============================================
-- 任务列表权限 -> 任务管理主菜单
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 任务详情权限 -> 任务管理查看按钮
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@detail_permission_id, @view_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 任务取消权限 -> 任务管理取消按钮
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@cancel_permission_id, @cancel_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 6. 插入权限-接口关联数据
-- ============================================
-- 任务列表权限 -> GET /api/v1/tasks接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 任务详情权限 -> GET /api/v1/tasks/detail接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@detail_permission_id, @detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 任务取消权限 -> POST /api/v1/tasks/cancel接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@cancel_permission_id, @cancel_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

