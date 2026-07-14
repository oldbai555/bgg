-- ============================================
-- 字典SQL增量脚本
-- 模块：性能日志/埋点统计
-- 创建时间：2026-07-14
-- 说明：admin-frontend MetricStats.vue「业务模块」筛选下拉此前硬编码四个 <el-option>，
--       未走字典驱动（违反 20-frontend.md 下拉必须来自字典的强制规则）。
--       value 直接对应后端 internal/consts/blog.go 里的 MetricModuleXxx 常量字符串，
--       与 sdk_http_method 字典（value 同样是 'GET'/'POST' 字符串而非数字）是同一种模式，
--       不受"字典 value 从 1 开始"规则约束（该规则只管数字型枚举）。
-- ============================================

INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('性能日志业务模块', 'metric_module', '性能日志/埋点统计的业务模块标识，对应后端 MetricModuleXxx 常量', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @metric_module_type_id = (
  SELECT `id` FROM `admin_dict_type`
  WHERE `code` = 'metric_module' AND `deleted_at` = 0
  LIMIT 1
);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@metric_module_type_id, '博客文章列表', 'blog_article_list', 1, 1, '对应 MetricModuleBlogArticleList', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@metric_module_type_id, '博客文章详情', 'blog_article_detail', 2, 1, '对应 MetricModuleBlogArticleDetail', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@metric_module_type_id, '视频列表', 'video_list', 3, 1, '对应 MetricModuleVideoList', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@metric_module_type_id, '视频详情', 'video_detail', 4, 1, '对应 MetricModuleVideoDetail', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
