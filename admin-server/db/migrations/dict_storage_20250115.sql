-- ============================================
-- 字典SQL增量脚本
-- 模块：文件存储
-- 创建时间：2025-01-15
-- 说明：本地存储配置字典，用于配置文件存储的baseURL，替代配置文件中的BaseURL
-- ============================================

-- 1. 插入字典类型
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('本地存储', 'storage_base_url', '本地存储配置，用于配置文件存储的baseURL', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 获取字典类型ID（通过code获取）
SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'storage_base_url' AND `deleted_at` = 0 LIMIT 1);

-- 3. 插入字典项（使用自动增长ID，通过code获取的type_id）
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

