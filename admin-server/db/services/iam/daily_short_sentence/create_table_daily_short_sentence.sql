-- iam/daily_short_sentence 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `daily_short_sentence` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `type` INT NOT NULL DEFAULT 1 COMMENT '类型：1普通，2文学',
  `content` TEXT NOT NULL COMMENT '短句内容',
  `img` TEXT COMMENT '图片URL',
  `literature_author` VARCHAR(255) COMMENT '作者',
  `convert_img` TEXT COMMENT '转换后的图片URL',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_daily_short_sentence_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日短句表';
