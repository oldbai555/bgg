-- iam/department 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 根部门
INSERT INTO `admin_department` (`id`, `parent_id`, `name`, `order_num`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (1, 0, '总部', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;
