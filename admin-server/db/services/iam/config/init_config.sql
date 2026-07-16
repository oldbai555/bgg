-- iam/config 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 系统配置初始化数据
INSERT INTO `admin_config` (`id`, `group`, `key`, `value`, `type`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, 'system', 'system:app_name', '"后台管理系统"', 'string', '应用名称', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'system', 'system:app_logo', '"/static/logo.png"', 'string', '应用Logo路径', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, 'system', 'system:app_version', '"1.0.0"', 'string', '应用版本', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 'system', 'system:timeout', '300', 'number', '会话超时时间（秒）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 'theme', 'theme:primary_color', '"#409EFF"', 'string', '主题主色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, 'theme', 'theme:sidebar_width', '200', 'number', '侧边栏宽度（px）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 'upload', 'upload:max_size', '10485760', 'number', '最大上传文件大小（字节，默认10MB）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, 'upload', 'upload:allowed_types', '["jpg","jpeg","png","gif","pdf","doc","docx","xls","xlsx"]', 'json', '允许上传的文件类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `value`=VALUES(`value`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
