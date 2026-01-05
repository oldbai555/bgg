CREATE TABLE IF NOT EXISTS `video` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(255) NOT NULL COMMENT '视频名称',
  `cover` VARCHAR(512) COMMENT '视频封面URL',
  `duration` INT NOT NULL DEFAULT 0 COMMENT '视频时长（秒）',
  `play_url` VARCHAR(512) NOT NULL COMMENT '播放链接',
  `description` TEXT COMMENT '视频描述',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_video_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频列表管理表';

