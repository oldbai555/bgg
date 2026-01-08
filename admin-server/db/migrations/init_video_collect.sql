-- 视频采集模块初始化 SQL
-- 功能组: video_collect
-- 功能名称: 视频采集相关接口（仅记录操作日志，不需要权限控制）
-- 说明：
-- 1. 采集视频已整合到统一的视频管理页面（/video/list），使用 type 筛选区分
-- 2. 采集接口本身不需要权限控制（仅记录操作日志）
-- 3. 后台管理通过现有的视频管理页面进行，使用现有的 video:* 权限

-- ============================================
-- 插入接口数据（仅用于记录操作日志）
-- ============================================
-- 视频采集接口（仅记录操作日志，不需要权限控制）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '视频采集',
    'POST',
    '/api/v1/videos/collect',
    '采集视频接口（仅记录操作日志）',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- m3u8代理接口（无中间件，不需要权限控制）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    'm3u8代理',
    'GET',
    '/api/v1/m3u8/proxy',
    'm3u8代理服务（无权限控制）',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- m3u8代理OPTIONS接口（CORS预检）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    'm3u8代理OPTIONS',
    'OPTIONS',
    '/api/v1/m3u8/proxy',
    'm3u8代理CORS预检请求（无权限控制）',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 公开视频列表接口（无中间件，不需要权限控制）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '公开视频列表',
    'GET',
    '/api/v1/public/videos/list',
    '公开视频列表（无权限控制）',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 公开视频详情接口（无中间件，不需要权限控制）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '公开视频详情',
    'GET',
    '/api/v1/public/videos/info/:id',
    '公开视频详情（无权限控制）',
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- ============================================
-- 说明：
-- 1. 采集接口（/api/v1/videos/collect）仅记录操作日志，不需要权限控制
-- 2. m3u8代理接口（/api/v1/m3u8/proxy）无中间件，不需要权限控制
-- 3. 公开视频接口（/api/v1/public/videos/*）无中间件，不需要权限控制
-- 4. 后台管理通过现有的视频管理页面（/video/list）进行，使用现有的 video:* 权限
-- 5. 采集视频通过 type=2 筛选条件在统一管理页面中查看和管理
-- ============================================

