-- content/video 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    0,
    '影视资源',
    '/admin/video',
    '',
    'ele-VideoPlay',
    1,
    20,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `path`=VALUES(`path`),
  `icon`=VALUES(`icon`),
  `order_num`=VALUES(`order_num`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 父菜单固定是本文件上面刚 UPSERT 的「影视资源」根目录，不像 daily_short_sentence/
-- metric 那类脚手架模块需要回退到临时目录（id=9）——回退到 9 是复制脚手架模板遗留的
-- 错误写法，一旦 /admin/video 查找意外落空会把视频子菜单错挂到临时目录下。
SET @parent_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/video' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '视频列表管理',
    '/admin/video/list',
    'content/VideoList',
    'ele-Document',
    2,
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/video/list' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@main_menu_id, '视频列表管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '视频列表管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '视频列表管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 新增按钮' AND `deleted_at`=0 LIMIT 1);
SET @update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 编辑按钮' AND `deleted_at`=0 LIMIT 1);
SET @delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 删除按钮' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频列表管理列表', 'video:list', '查看视频列表管理列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理新增', 'video:create', '新增视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理编辑', 'video:update', '编辑视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理删除', 'video:delete', '删除视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:list' AND `deleted_at`=0 LIMIT 1);
SET @create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:create' AND `deleted_at`=0 LIMIT 1);
SET @update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:update' AND `deleted_at`=0 LIMIT 1);
SET @delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频列表管理列表', 'GET', '/api/v1/videos', '获取视频列表管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理新增', 'POST', '/api/v1/videos', '新增视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理编辑', 'PUT', '/api/v1/videos', '编辑视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理删除', 'DELETE', '/api/v1/videos', '删除视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频采集', 'POST', '/api/v1/videos/collect', '采集视频接口（仅记录操作日志）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('m3u8代理', 'GET', '/api/v1/m3u8/proxy', 'm3u8代理服务（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('m3u8代理OPTIONS', 'OPTIONS', '/api/v1/m3u8/proxy', 'm3u8代理CORS预检请求（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公开视频列表', 'GET', '/api/v1/public/videos/list', '公开视频列表（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公开视频详情', 'GET', '/api/v1/public/videos/info', '公开视频详情（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='GET' AND `deleted_at`=0 LIMIT 1);
SET @create_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='POST' AND `deleted_at`=0 LIMIT 1);
SET @update_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='PUT' AND `deleted_at`=0 LIMIT 1);
SET @delete_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='DELETE' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '视频播放器',
    '/admin/video/player',
    'content/VideoPlayer',
    'ele-VideoPlay',
    2,
    1,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

