-- iam/department 增量：新增"飞书待分配"部门，供飞书扫码登录首次自动建号时默认分配。
-- admin_department 没有 name 唯一键（历史遗留，见 docs/changelog/archive-backend.md 第 23 节同类问题
-- 记录），沿用 init_department.sql 已有的做法——显式指定 id，让 ON DUPLICATE KEY UPDATE
-- 命中主键才能真正幂等，不能像这张表理论上该有唯一键那样直接套用 name 判重。

INSERT INTO `admin_department` (`id`, `parent_id`, `name`, `order_num`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (2, 0, '飞书待分配', 99, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at` = 0, `name` = '飞书待分配';
