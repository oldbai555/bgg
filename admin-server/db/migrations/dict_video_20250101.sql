-- ============================================
-- 字典SQL增量脚本
-- 模块：视频管理
-- 创建时间：2025-01-01
-- 说明：视频代理地址字典，用于配置m3u8视频代理服务器地址
-- ============================================

-- 1. 插入字典类型（使用自动增长ID）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频代理地址', 'video_proxy_url', '视频代理服务器地址配置，用于代理m3u8等视频流', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 获取字典类型ID（通过code获取）
SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'video_proxy_url' AND `deleted_at` = 0 LIMIT 1);

-- 3. 插入字典项（使用自动增长ID，通过code获取的type_id）
INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认代理', 'http://localhost:8888/api/v1/videos/proxy', 1, 1, '默认视频代理服务器地址', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

