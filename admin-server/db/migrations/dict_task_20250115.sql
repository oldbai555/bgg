-- ============================================
-- 字典SQL增量脚本
-- 模块：异步任务管理
-- 创建时间：2025-01-15
-- 说明：异步任务相关字典，包括任务类型、执行类型、状态、配置等
-- ============================================

-- 1. 插入字典类型：任务类型（task_type）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('任务类型', 'task_type', '异步任务类型字典，如：异步导出Excel、定时发送邮件通知等', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 获取字典类型ID（通过code获取）
SET @task_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'task_type' AND `deleted_at` = 0 LIMIT 1);

-- 3. 插入字典项：任务类型（task_type）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@task_type_id, '异步导出Excel', '1', 1, 1, '异步导出Excel文件', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 4. 插入字典类型：任务执行类型（task_execution_type）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('任务执行类型', 'task_execution_type', '任务执行类型字典：同步、异步', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 5. 获取字典类型ID（通过code获取）
SET @task_execution_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'task_execution_type' AND `deleted_at` = 0 LIMIT 1);

-- 6. 插入字典项：任务执行类型（task_execution_type）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@task_execution_type_id, '同步', '1', 1, 1, '同步执行任务', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@task_execution_type_id, '异步', '2', 2, 1, '异步执行任务', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 7. 插入字典类型：任务状态（task_status）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('任务状态', 'task_status', '任务状态字典：未开始、进行中、已完成、失败', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8. 获取字典类型ID（通过code获取）
SET @task_status_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'task_status' AND `deleted_at` = 0 LIMIT 1);

-- 9. 插入字典项：任务状态（task_status）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@task_status_id, '未开始', '1', 1, 1, '任务已创建，未开始执行', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@task_status_id, '进行中', '2', 2, 1, '任务正在执行中', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@task_status_id, '已完成', '3', 3, 1, '任务执行成功', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@task_status_id, '失败', '4', 4, 1, '任务执行失败', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 10. 插入字典类型：任务配置（task_config）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('任务配置', 'task_config', '任务系统配置字典，如：最近任务数量等', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 11. 获取字典类型ID（通过code获取）
SET @task_config_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'task_config' AND `deleted_at` = 0 LIMIT 1);

-- 12. 插入字典项：任务配置（task_config）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@task_config_id, '最近任务数量', '10', 1, 1, '浮动球显示最近N个任务（默认10个）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 13. 插入字典类型：通知来源类型（notification_source_type）- 如果不存在则添加 "task" 选项
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('通知来源类型', 'notification_source_type', '消息通知来源类型字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 14. 获取字典类型ID（通过code获取）
SET @notification_source_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notification_source_type' AND `deleted_at` = 0 LIMIT 1);

-- 15. 插入字典项：通知来源类型 - 任务（task）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
SELECT
  @notification_source_type_id,
  '任务',
  'task',
  COALESCE((SELECT MAX(`sort`) FROM `admin_dict_item` WHERE `type_id` = @notification_source_type_id AND `deleted_at` = 0), 0) + 1,
  1,
  '任务通知来源类型',
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_dict_item`
  WHERE `type_id` = @notification_source_type_id
  AND `value` = 'task'
  AND `deleted_at` = 0
)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

