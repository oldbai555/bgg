-- ============================================
-- 视频表字段扩展迁移脚本
-- 创建时间：2025-01-20
-- 说明：为现有video表添加采集视频相关字段，使用type字段区分来源
-- ============================================

-- 为现有video表添加新字段
ALTER TABLE `video` 
  ADD COLUMN `uuid` VARCHAR(36) COMMENT '唯一标识（采集视频使用，可为空）' AFTER `id`,
  ADD COLUMN `god_num` VARCHAR(127) COMMENT '番号（采集视频使用）' AFTER `name`,
  ADD COLUMN `xlzz_urls` JSON COMMENT '磁力链接数组（采集视频使用）' AFTER `play_url`,
  ADD COLUMN `type` TINYINT DEFAULT 1 COMMENT '来源类型：1=手动添加，2=采集（从字典获取）' AFTER `description`;

-- 添加唯一索引（如果不存在）
-- 注意：MySQL的唯一索引允许多个NULL值，所以uuid字段可以为空
ALTER TABLE `video` 
  ADD UNIQUE KEY `uk_uuid` (`uuid`);

-- 添加type字段索引（如果不存在）
ALTER TABLE `video` 
  ADD KEY `idx_type` (`type`);

-- 迁移现有video_sts表数据（如果有）
-- 注意：如果video_sts表不存在，此语句会报错，可以忽略
INSERT INTO `video` (`uuid`, `name`, `god_num`, `play_url`, `xlzz_urls`, `type`, `created_at`, `updated_at`, `deleted_at`)
SELECT `uuid`, `name`, `god_num`, `player_url`, `xlzz_urls`, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
FROM `video_sts`
WHERE NOT EXISTS (SELECT 1 FROM `video` WHERE `video`.`uuid` = `video_sts`.`uuid` AND `video_sts`.`uuid` IS NOT NULL);

