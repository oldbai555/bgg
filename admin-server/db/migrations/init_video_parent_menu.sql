-- 创建"影视资源"父目录
-- 用于挂载视频列表管理和视频播放器菜单

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    0, -- 根目录
    '影视资源',
    '/video',
    '',
    'ele-VideoPlay',
    1, -- 类型：1 目录
    20, -- 排序值（在系统管理之后）
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
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

