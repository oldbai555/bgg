-- content/video 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `video` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `uuid` VARCHAR(36) COMMENT '唯一标识（采集视频使用，可为空）',
  `name` VARCHAR(255) NOT NULL COMMENT '视频名称',
  `god_num` VARCHAR(127) COMMENT '番号（采集视频使用）',
  `cover` VARCHAR(512) COMMENT '视频封面URL',
  `duration` INT NOT NULL DEFAULT 0 COMMENT '视频时长（秒）',
  `play_url` VARCHAR(512) NOT NULL COMMENT '播放链接',
  `xlzz_urls` JSON COMMENT '磁力链接数组（采集视频使用）',
  `description` TEXT COMMENT '视频描述',
  `type` TINYINT DEFAULT 1 COMMENT '来源类型：1=手动添加，2=采集（从字典获取）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uuid` (`uuid`),
  KEY `idx_type` (`type`),
  KEY `idx_video_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频列表管理表';
