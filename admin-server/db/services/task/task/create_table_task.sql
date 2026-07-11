-- ============================================
-- 异步任务表结构SQL
-- 模块：异步任务管理
-- 创建时间：2025-01-15
-- 说明：异步任务表，用于存储任务信息、参数、结果等
-- ============================================

CREATE TABLE IF NOT EXISTS `admin_task` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(255) NOT NULL COMMENT '任务名称',
  `type` INT NOT NULL COMMENT '任务类型（字典：task_type，如：1=异步导出Excel，2=定时发送邮件通知）',
  `execution_type` INT NOT NULL COMMENT '执行类型（字典：task_execution_type，1=同步，2=异步）',
  `status` INT NOT NULL DEFAULT 1 COMMENT '任务状态（字典：task_status，1=未开始，2=进行中，3=已完成，4=失败）',
  `params` TEXT COMMENT '任务参数（JSON格式，存储任务执行所需的参数）',
  `result` TEXT COMMENT '任务结果（JSON格式，存储任务执行结果）',
  `error_message` VARCHAR(1000) DEFAULT '' COMMENT '错误信息（任务失败时存储）',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建用户ID',
  `scheduled_at` BIGINT NOT NULL DEFAULT 0 COMMENT '计划执行时间（秒级时间戳，0表示立即执行）',
  `started_at` BIGINT NOT NULL DEFAULT 0 COMMENT '开始执行时间（秒级时间戳）',
  `finished_at` BIGINT NOT NULL DEFAULT 0 COMMENT '完成时间（秒级时间戳）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_admin_task_type` (`type`),
  KEY `idx_admin_task_status` (`status`),
  KEY `idx_admin_task_user_id` (`user_id`),
  KEY `idx_admin_task_scheduled_at` (`scheduled_at`),
  KEY `idx_admin_task_created_at` (`created_at`),
  KEY `idx_admin_task_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='异步任务表';

