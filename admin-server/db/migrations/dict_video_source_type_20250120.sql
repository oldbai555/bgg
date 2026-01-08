-- ============================================
-- 字典SQL增量脚本
-- 模块：视频管理
-- 创建时间：2025-01-20
-- 说明：视频来源类型字典，用于区分手动添加和采集的视频
-- ============================================

-- 1. 插入字典类型（使用自动增长ID）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频来源类型', 'video_source_type', '视频来源类型字典，用于区分手动添加和采集的视频', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 获取字典类型ID（通过code获取）
SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'video_source_type' AND `deleted_at` = 0 LIMIT 1);

-- 3. 插入字典项（使用自动增长ID，通过code获取的type_id，枚举值从1开始）
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

