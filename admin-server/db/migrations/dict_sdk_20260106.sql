-- ============================================
-- 字典SQL增量脚本
-- 模块：SDK 管理
-- 创建时间：2026-01-06
-- 说明：SDK 状态、默认限频配置
-- ============================================

-- 1. SDK 状态字典类型与项
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

-- 2. SDK 默认限频字典（接口默认值兜底，单位：次/分钟）
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

-- 3. SDK HTTP 方法
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

