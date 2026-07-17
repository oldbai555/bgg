-- iam/metric 初始化数据
-- 功能组: metric

-- ============================================
-- 1. metric_daily_stats 没有独立菜单，也没有对应的真实接口
-- ============================================
-- 本文件历史上曾为一个名为 metric_daily_stats 的脚手架示例（views/temp/MetricList.vue，
-- 已在 admin-frontend Phase 1 Week 2 死代码清理中删除，见
-- admin-frontend/docs/07-cleanup-and-tooling.md §1）种过 metric:list/create/update/delete
-- 权限 + /api/v1/metrics[/:id] 接口数据，但 admin-server/api/admin.api 里从未真正存在过
-- 这套 CRUD 路由（metric 模块实际只有 /metrics/report、/metrics/stats 两个接口，
-- 见 MetricReport/MetricStats 服务块），代码库里也找不到任何 MetricCreate/MetricUpdate/
-- MetricDelete 的 Go 符号引用——这是纯粹的孤儿种子数据，本次连同孤儿菜单一起删除，
-- 不再保留。

-- ============================================
-- 2. 数据统计（PV/UV/VV/IP）菜单 + 权限 + 接口
-- ============================================
-- 原本单独放在 init_metric_stats.sql，但 db/services/init-dev-db.sh 的 run_module()
-- 只按 init_<module>.sql 命名约定自动执行（本模块目录名是 metric，不是
-- metric_stats），导致 init_metric_stats.sql 从未被自动执行过、「数据统计」菜单在
-- 全新库里根本不会出现，本次并入本文件修复这个遗漏。

SET @system_dir_id = (SELECT `id` FROM `admin_menu` WHERE `id` = 2 AND `deleted_at` = 0 LIMIT 1);

-- 数据统计主菜单（放在系统管理下）；admin_menu 没有唯一键，ON DUPLICATE KEY UPDATE 对它
-- 不会触发，改用 INSERT ... SELECT ... WHERE NOT EXISTS 保证幂等（同 content/blog/init_blog.sql）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @system_dir_id, '数据统计', '/admin/system/metric-stats', 'monitoring/MetricStats', 'ele-DataAnalysis', 2, 21, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `path` = '/admin/system/metric-stats' AND `deleted_at` = 0);

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/system/metric-stats' AND `deleted_at` = 0 LIMIT 1);

-- 数据统计查询权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('数据统计', 'metric:stats', '查看PV/UV/VV/IP统计数据', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @stats_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'metric:stats' AND `deleted_at` = 0 LIMIT 1);

-- 数据统计接口（接口可能已通过路由同步自动创建，这里使用 ON DUPLICATE KEY UPDATE）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('数据统计', 'GET', '/api/v1/metrics/stats', '获取PV/UV/VV/IP统计数据', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @stats_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/metrics/stats' AND `deleted_at` = 0 LIMIT 1);

-- 数据统计权限 -> 数据统计主菜单
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@stats_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 数据统计权限 -> GET /api/v1/metrics/stats接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@stats_permission_id, @stats_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();
