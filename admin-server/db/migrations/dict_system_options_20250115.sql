-- ============================================
-- 字典SQL增量脚本
-- 模块：系统选项字典
-- 创建时间：2025-01-15
-- 说明：系统各模块el-select选项字典化，统一管理所有下拉选项
-- ============================================

-- 1. 性能日志慢查询状态
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('性能日志慢查询状态', 'performance_log_slow_status', '性能日志慢查询状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'performance_log_slow_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 枚举从 1 开始，0 预留为「全部/不筛选」
  (@dict_type_id, 'Normal', '2', 1, 1, '正常查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, 'Slow', '1', 2, 1, '慢查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 菜单类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('菜单类型', 'menu_type', '菜单类型选项：目录、菜单、按钮', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'menu_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '目录', '1', 1, 1, '目录类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '菜单', '2', 2, 1, '菜单类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '按钮', '3', 3, 1, '按钮类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 3. 已读状态
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('已读状态', 'read_status', '已读状态选项：未读、已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'read_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 枚举从 1 开始，0 预留为「全部/不筛选」
  (@dict_type_id, '未读', '1', 1, 1, '未读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '已读', '2', 2, 1, '已读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 4. 操作类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('操作类型', 'operation_type', '操作日志操作类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'operation_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '创建', 'create', 1, 1, '创建操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '更新', 'update', 2, 1, '更新操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '删除', 'delete', 3, 1, '删除操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '查询', 'query', 4, 1, '查询操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '导出', 'export', 5, 1, '导出操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 5. HTTP请求方法
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('HTTP请求方法', 'http_method', 'HTTP请求方法选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'http_method' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, 'GET', 'GET', 1, 1, 'GET请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, 'POST', 'POST', 2, 1, 'POST请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, 'PUT', 'PUT', 3, 1, 'PUT请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, 'DELETE', 'DELETE', 4, 1, 'DELETE请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 6. 公告类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('公告类型', 'notice_type', '公告类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '普通公告', '1', 1, 1, '普通公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '重要公告', '2', 2, 1, '重要公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '紧急公告', '3', 3, 1, '紧急公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 7. 公告状态
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('公告状态', 'notice_status', '公告状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '草稿', '1', 1, 1, '草稿状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '已发布', '2', 2, 1, '已发布状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8. 登录状态
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('登录状态', 'login_status', '登录状态选项：成功、失败', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'login_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 枚举从 1 开始，0 预留为「全部/不筛选」
  (@dict_type_id, '失败', '2', 1, 1, '登录失败', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '成功', '1', 2, 1, '登录成功', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 9. 审计类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('审计类型', 'audit_type', '审计日志审计类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'audit_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '权限分配', 'permission_assign', 1, 1, '权限分配审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '角色变更', 'role_change', 2, 1, '角色变更审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '配置修改', 'config_modify', 3, 1, '配置修改审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '数据删除', 'data_delete', 4, 1, '数据删除审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 10. 消息类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('消息类型', 'chat_message_type', '聊天消息类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'chat_message_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '文本', '1', 1, 1, '文本消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '图片', '2', 2, 1, '图片消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '文件', '3', 3, 1, '文件消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 11. 短句类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('短句类型', 'daily_short_sentence_type', '每日短句类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'daily_short_sentence_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '普通', '1', 1, 1, '普通短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '文学', '2', 2, 1, '文学短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

