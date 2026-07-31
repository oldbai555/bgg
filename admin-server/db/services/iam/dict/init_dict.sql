-- iam/dict 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 数据字典类型初始化数据
INSERT INTO `admin_dict_type` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, '用户状态', 'user_status', '用户账号状态字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, '性别', 'gender', '性别字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, '是否', 'yes_no', '是否字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, '文件存储类型', 'file_storage_type', '文件存储类型字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, '聊天配置', 'chat_config', '在线聊天相关配置字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, '消息来源类型', 'notification_source_type', '消息通知的来源类型字典（chat、notice、system等）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 数据字典项初始化数据
INSERT INTO `admin_dict_item` (`id`, `type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 用户状态字典项（枚举从 1 开始，0 预留为「全部/不筛选」）
  (1, 1, '启用', '1', 1, 1, '用户账号启用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 1, '禁用', '2', 2, 1, '用户账号禁用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 性别字典项（枚举从 1 开始，0 不再使用）
  (3, 2, '男', '1', 1, 1, '男性', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 2, '女', '2', 2, 1, '女性', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 2, '未知', '3', 3, 1, '未知性别', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 是否字典项（枚举从 1 开始，0 不再使用）
  (6, 3, '是', '1', 1, 1, '是', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 3, '否', '2', 2, 1, '否', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 文件存储类型字典项
  (8, 4, '本地存储', 'local', 1, 1, '本地文件系统存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, 4, 'OSS存储', 'oss', 2, 1, '阿里云OSS存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (10, 4, 'S3存储', 's3', 3, 1, 'AWS S3存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 聊天配置字典项
  (11, 5, '聊天窗口消息数量', '30', 1, 1, '每个聊天窗口显示的最新消息数量', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, 5, '在线聊天页面路径', '/admin/chatroom/chat', 2, 1, '在线聊天页面的前端路由路径，用于消息通知跳转', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, 5, 'Emoji每行显示数量', '8', 3, 1, 'Emoji表情选择器每行显示的表情数量（x）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, 5, 'Emoji显示行数', '3', 4, 1, 'Emoji表情选择器显示的行数（y）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 消息来源类型字典项
  (13, 6, '在线聊天', 'chat', 1, 1, '在线聊天消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, 6, '系统公告', 'notice', 2, 1, '系统公告消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, 6, '系统通知', 'system', 3, 1, '系统通知消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- ============================================
-- 8. 字典数据（SDK/视频/存储/WebSocket/日志等模块的字典类型与字典项）
--
-- 这部分字典的业务归属模块不是 iam（如 sdk_status 属于 sdk 模块、video_source_type
-- 属于 content/video 模块），但 admin_dict_type/admin_dict_item 两张物理表本身是
-- iam 拥有的全局共享表，所以按"物理表归属"落在 iam/dict 目录下，不是按"业务归属"
-- 拆到各自模块目录。这些字典条目在项目最初上线时就已经和上面的基础字典一起写在
-- db/data.sql 里（不是后来以增量迁移形式追加的），因此随首次部署的 init_dict.sql
-- 一起写入是准确的，不属于"字典 SQL 插入 data.sql 破坏增量边界"的反例——该反例针对
-- 的是"已上线后才新增的字典改动 data.sql"，不是"首次部署阶段内的字典种子数据"。
-- 后续任何模块新增字典，一律走 db/services/<service>/<module>/migrations/
-- dict_{module}_YYYYMMDD.sql，不再追加进这里（对照 content/blog 的
-- migrations/dict_blog_20260114.sql 先例）。
-- ============================================

-- 8.1 SDK 状态字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK状态', 'sdk_status', 'SDK Key 状态（启用/禁用）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_status_type_id, '启用', '1', 1, 1, '可用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_status_type_id, '禁用', '2', 2, 1, '停用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.2 SDK 默认限频字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK默认限频', 'sdk_rate_limit_default', 'SDK 接口默认限频（次/分钟），可被接口自定义值覆盖', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_rate_limit_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_rate_limit_default' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_rate_limit_type_id, '默认60次/分钟', '60', 1, 1, '默认限频上限，单位：次/分钟', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.3 SDK HTTP 方法
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK HTTP 方法', 'sdk_http_method', 'SDK 接口支持的 HTTP 方法', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_http_method_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_http_method' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_http_method_type_id, 'GET', 'GET', 1, 1, 'HTTP GET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'POST', 'POST', 2, 1, 'HTTP POST', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'PUT', 'PUT', 3, 1, 'HTTP PUT', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'DELETE', 'DELETE', 4, 1, 'HTTP DELETE', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.4 本地存储配置字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('本地存储', 'storage_base_url', '本地存储配置，用于配置文件存储的baseURL', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'storage_base_url' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认存储地址', 'https://oldbai.top/oss', 1, 1, '文件存储的baseURL，用于生成文件完整访问路径', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.5 系统选项字典（性能日志慢查询状态、菜单类型、已读状态、操作类型、HTTP请求方法、公告类型、公告状态、登录状态、审计类型、消息类型、短句类型）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('性能日志慢查询状态', 'performance_log_slow_status', '性能日志慢查询状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('菜单类型', 'menu_type', '菜单类型选项：目录、菜单、按钮', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('已读状态', 'read_status', '已读状态选项：未读、已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作类型', 'operation_type', '操作日志操作类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('HTTP请求方法', 'http_method', 'HTTP请求方法选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告类型', 'notice_type', '公告类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告状态', 'notice_status', '公告状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录状态', 'login_status', '登录状态选项：成功、失败', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计类型', 'audit_type', '审计日志审计类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息类型', 'chat_message_type', '聊天消息类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('短句类型', 'daily_short_sentence_type', '每日短句类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @performance_log_slow_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'performance_log_slow_status' AND `deleted_at` = 0 LIMIT 1);
SET @menu_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'menu_type' AND `deleted_at` = 0 LIMIT 1);
SET @read_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'read_status' AND `deleted_at` = 0 LIMIT 1);
SET @operation_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'operation_type' AND `deleted_at` = 0 LIMIT 1);
SET @http_method_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'http_method' AND `deleted_at` = 0 LIMIT 1);
SET @notice_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_type' AND `deleted_at` = 0 LIMIT 1);
SET @notice_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_status' AND `deleted_at` = 0 LIMIT 1);
SET @login_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'login_status' AND `deleted_at` = 0 LIMIT 1);
SET @audit_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'audit_type' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'chat_message_type' AND `deleted_at` = 0 LIMIT 1);
SET @daily_short_sentence_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'daily_short_sentence_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@performance_log_slow_status_type_id, 'Normal', '2', 1, 1, '正常查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@performance_log_slow_status_type_id, 'Slow', '1', 2, 1, '慢查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '目录', '1', 1, 1, '目录类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '菜单', '2', 2, 1, '菜单类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '按钮', '3', 3, 1, '按钮类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@read_status_type_id, '未读', '1', 1, 1, '未读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@read_status_type_id, '已读', '2', 2, 1, '已读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '创建', 'create', 1, 1, '创建操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '更新', 'update', 2, 1, '更新操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '删除', 'delete', 3, 1, '删除操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '查询', 'query', 4, 1, '查询操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '导出', 'export', 5, 1, '导出操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'GET', 'GET', 1, 1, 'GET请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'POST', 'POST', 2, 1, 'POST请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'PUT', 'PUT', 3, 1, 'PUT请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'DELETE', 'DELETE', 4, 1, 'DELETE请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '普通公告', '1', 1, 1, '普通公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '重要公告', '2', 2, 1, '重要公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '紧急公告', '3', 3, 1, '紧急公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_status_type_id, '草稿', '1', 1, 1, '草稿状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_status_type_id, '已发布', '2', 2, 1, '已发布状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@login_status_type_id, '失败', '2', 1, 1, '登录失败', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@login_status_type_id, '成功', '1', 2, 1, '登录成功', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '权限分配', 'permission_assign', 1, 1, '权限分配审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '角色变更', 'role_change', 2, 1, '角色变更审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '配置修改', 'config_modify', 3, 1, '配置修改审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '数据删除', 'data_delete', 4, 1, '数据删除审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '文本', '1', 1, 1, '文本消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '图片', '2', 2, 1, '图片消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '文件', '3', 3, 1, '文件消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@daily_short_sentence_type_type_id, '普通', '1', 1, 1, '普通短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@daily_short_sentence_type_type_id, '文学', '2', 2, 1, '文学短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.7 视频来源类型字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频来源类型', 'video_source_type', '视频来源类型字典，用于区分手动添加和采集的视频', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'video_source_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '手动添加', '1', 1, 1, '手动添加的视频', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '采集', '2', 2, 1, '通过采集接口添加的视频', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.8 WebSocket连接配置字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('WebSocket连接', 'websocket_base_url', 'WebSocket连接配置，用于配置WebSocket的baseURL', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'websocket_base_url' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认WebSocket地址', 'oldbai.top/ws', 1, 1, 'WebSocket连接的baseURL，生产环境使用 wss://oldbai.top/ws，开发环境使用 ws://localhost:20000', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

