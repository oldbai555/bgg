CREATE TABLE IF NOT EXISTS `metric_daily_stats` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `module` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务模块标识，如 blog_article_list/blog_article_detail/video_list/video_detail',
  `biz_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '业务ID（文章ID、视频ID等；0表示列表页）',
  `day` CHAR(8) NOT NULL DEFAULT '' COMMENT '统计日期：YYYYMMDD',
  `pv` BIGINT NOT NULL DEFAULT 0 COMMENT '页面访问量（Page View）',
  `uv` BIGINT NOT NULL DEFAULT 0 COMMENT '独立访客数（Unique Visitor，基于 IP + User-Agent 去重）',
  `vv` BIGINT NOT NULL DEFAULT 0 COMMENT '访问次数/播放次数（Visit/View）',
  `ip` BIGINT NOT NULL DEFAULT 0 COMMENT '独立IP数量（基于 IP 去重）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_metric_stats_module_biz_day` (`module`, `biz_id`, `day`),
  KEY `idx_metric_stats_day` (`day`),
  KEY `idx_metric_stats_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='PV/UV/VV/IP 日统计表';

