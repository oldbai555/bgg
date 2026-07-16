-- 把 task/sdk/chat/content 四个拆分服务的表，全部建在现有的 oldbai 库里（sqlpub 免费版只有 1 个库）
-- 由用户手动执行： mysql -h mysql4.sqlpub.com -P 3309 -u oldbai -p oldbai < init-tables-into-oldbai.sql
-- 幂等：全部 CREATE TABLE IF NOT EXISTS，重复执行安全；表名与 iam 现有表不冲突

-- task
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


-- sdk
-- sdk/sdk 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `sdk_key` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',
  `api_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '对外 API Key（唯一）',
  `api_secret` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '对外 API Secret（唯一）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `expire_at` BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间(秒级时间戳)，0 表示长期有效',
  `ip_whitelist` TEXT NOT NULL COMMENT 'IP 白名单，逗号或换行分隔，空表示不限制',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_key_api_key` (`api_key`),
  UNIQUE KEY `uk_sdk_key_api_secret` (`api_secret`),
  KEY `idx_sdk_key_status` (`status`),
  KEY `idx_sdk_key_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK Key';

CREATE TABLE IF NOT EXISTS `sdk_interface` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口名称',
  `api_code` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口编码（唯一，便于绑定与校验）',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径（对外 SDK 路由）',
  `method` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'HTTP 方法',
  `rate_limit_default` INT NOT NULL DEFAULT 0 COMMENT '默认限频上限（0 表示不限，实际逻辑读取字典兜底）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_interface_code` (`api_code`),
  KEY `idx_sdk_interface_status` (`status`),
  KEY `idx_sdk_interface_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK 对外接口';

CREATE TABLE IF NOT EXISTS `sdk_key_api` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `sdk_key_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK Key ID',
  `sdk_interface_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK 接口 ID',
  `custom_rate_limit` INT NOT NULL DEFAULT 0 COMMENT '自定义限频上限（0 表示使用接口默认值或字典默认值）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sdk_key_api` (`sdk_key_id`, `sdk_interface_id`),
  KEY `idx_sdk_key_api_key_id` (`sdk_key_id`),
  KEY `idx_sdk_key_api_interface_id` (`sdk_interface_id`),
  KEY `idx_sdk_key_api_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK Key 接口授权';

CREATE TABLE IF NOT EXISTS `sdk_call_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `sdk_key_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK Key ID',
  `sdk_interface_id` BIGINT UNSIGNED NOT NULL COMMENT 'SDK 接口 ID',
  `api_code` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口编码快照',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '请求路径',
  `method` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'HTTP 方法',
  `ip` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '调用 IP',
  `user_agent` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'User-Agent',
  `req_body` MEDIUMTEXT COMMENT '请求体（脱敏/截断后写入）',
  `resp_body` MEDIUMTEXT COMMENT '响应体（脱敏/截断后写入）',
  `resp_code` INT NOT NULL DEFAULT 0 COMMENT 'HTTP 状态码',
  `duration_ms` INT NOT NULL DEFAULT 0 COMMENT '耗时（毫秒）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_sdk_call_log_key` (`sdk_key_id`),
  KEY `idx_sdk_call_log_interface` (`sdk_interface_id`),
  KEY `idx_sdk_call_log_created_at` (`created_at`),
  KEY `idx_sdk_call_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SDK 调用日志';

-- chat
-- chat/chat 建表
-- 从 admin-server/db/tables.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

CREATE TABLE IF NOT EXISTS `chat` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '聊天名称（群组名称，私聊为空）',
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '聊天类型：1私聊，2群组',
  `avatar` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像URL（群组头像，私聊为空）',
  `description` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '描述（群组描述，私聊为空）',
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID（群组创建人，私聊为0）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_chat_type` (`type`),
  KEY `idx_chat_created_by` (`created_by`),
  KEY `idx_chat_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天表（私聊、群组）';

CREATE TABLE IF NOT EXISTS `chat_user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `chat_id` BIGINT UNSIGNED NOT NULL COMMENT '聊天 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
  `joined_at` BIGINT NOT NULL DEFAULT 0 COMMENT '加入时间(秒级时间戳)',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chat_user` (`chat_id`,`user_id`),
  KEY `idx_chat_user_chat_id` (`chat_id`),
  KEY `idx_chat_user_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天-用户关联表';

CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `chat_id` BIGINT UNSIGNED NOT NULL COMMENT '聊天 ID（关联chat表）',
  `from_user_id` BIGINT UNSIGNED NOT NULL COMMENT '发送用户 ID',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `message_type` TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型：1文本，2图片，3文件',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1已发送，2已读，3已撤回',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_chat_message_chat_id` (`chat_id`),
  KEY `idx_chat_message_from_user_id` (`from_user_id`),
  KEY `idx_chat_message_created_at` (`created_at`),
  KEY `idx_chat_message_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

-- content: blog
-- ============================================
-- 博客模块表结构 SQL
-- 模块：博客（标签、文章、审核）
-- 创建时间：2026-01-14
-- 说明：包含标签表、文章表、文章标签关联表、文章审核记录表
-- ============================================

-- 1. 文章标签表
CREATE TABLE IF NOT EXISTS `blog_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_tag_status）',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_tag_name` (`name`),
  KEY `idx_blog_tag_status` (`status`),
  KEY `idx_blog_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='博客标签表';

-- 2. 文章主表
CREATE TABLE IF NOT EXISTS `blog_article` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
  `content` LONGTEXT NOT NULL COMMENT '文章内容（Markdown 原文）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '文章状态（字典：blog_article_status）',
  `audit_status` TINYINT NOT NULL DEFAULT 1 COMMENT '审核状态（字典：blog_article_audit_status）',
  `cover` VARCHAR(512) DEFAULT '' COMMENT '封面图片 URL（为空时前端按标题首字生成占位封面）',
  `author_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '作者后台用户ID',
  `author_name` VARCHAR(64) DEFAULT '' COMMENT '作者姓名快照',
  `publish_time` BIGINT NOT NULL DEFAULT 0 COMMENT '上架时间(秒级时间戳,0表示未上架)',
  `summary` VARCHAR(512) DEFAULT '' COMMENT '文章摘要/简介（可选，公共列表优先使用）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_status` (`status`),
  KEY `idx_blog_article_audit_status` (`audit_status`),
  KEY `idx_blog_article_author_id` (`author_id`),
  KEY `idx_blog_article_publish_time` (`publish_time`),
  KEY `idx_blog_article_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='博客文章表';

-- 3. 文章-标签关联表（多对多）
CREATE TABLE IF NOT EXISTS `blog_article_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签 ID',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_tag_article_id` (`article_id`),
  KEY `idx_blog_article_tag_tag_id` (`tag_id`),
  KEY `idx_blog_article_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';

-- 4. 文章审核记录表
CREATE TABLE IF NOT EXISTS `blog_article_audit` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `audit_status` TINYINT NOT NULL COMMENT '本次审核结果（字典：blog_article_audit_status）',
  `audit_remark` VARCHAR(1000) DEFAULT '' COMMENT '审核意见',
  `auditor_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核人后台用户ID',
  `auditor_name` VARCHAR(64) DEFAULT '' COMMENT '审核人姓名快照',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_article_audit_article_id` (`article_id`),
  KEY `idx_blog_article_audit_auditor_id` (`auditor_id`),
  KEY `idx_blog_article_audit_created_at` (`created_at`),
  KEY `idx_blog_article_audit_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章审核记录表';


-- content: blog_extension
-- ============================================
-- 博客扩展功能表结构 SQL
-- 模块：博客扩展功能（友情链接、社交信息、文章置顶）
-- 创建时间：2026-01-16
-- 说明：包含友情链接表、社交信息表，以及文章表置顶字段扩展
-- ============================================

-- ============================================
-- 1. 友情链接表
-- ============================================
CREATE TABLE IF NOT EXISTS `blog_friend_link` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '链接名称（最大15个中文字符，通过字典配置）',
  `url` VARCHAR(255) NOT NULL COMMENT '目标链接（最大255字符，通过字典配置）',
  `remark` VARCHAR(512) DEFAULT '' COMMENT '备注（最大127个中文字符，通过字典配置）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_friend_link_status）',
  `order_num` INT NOT NULL DEFAULT 0 COMMENT '排序值（数字越小越靠前）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_friend_link_status` (`status`),
  KEY `idx_blog_friend_link_order_num` (`order_num`),
  KEY `idx_blog_friend_link_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='友情链接表';

-- ============================================
-- 2. 社交信息表
-- ============================================
CREATE TABLE IF NOT EXISTS `blog_social_info` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '社交平台名称（最大15个中文字符，通过字典配置）',
  `url` VARCHAR(255) NOT NULL COMMENT '目标链接（最大255字符，通过字典配置）',
  `remark` VARCHAR(512) DEFAULT '' COMMENT '备注（最大127个中文字符，通过字典配置）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=启用，2=禁用（字典：blog_social_info_status）',
  `order_num` INT NOT NULL DEFAULT 0 COMMENT '排序值（数字越小越靠前）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_blog_social_info_status` (`status`),
  KEY `idx_blog_social_info_order_num` (`order_num`),
  KEY `idx_blog_social_info_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社交信息表';

-- ============================================
-- 3. 文章表扩展（添加置顶字段）
-- ============================================
-- 检查字段是否已存在，如果不存在则添加
SET @dbname = DATABASE();
SET @tablename = 'blog_article';
SET @columnname = 'is_top';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname)
  ) > 0,
  'SELECT 1',
  CONCAT('ALTER TABLE `', @tablename, '` ADD COLUMN `', @columnname, '` TINYINT NOT NULL DEFAULT 0 COMMENT ''是否置顶：0=否，1=是'' AFTER `summary`, ADD KEY `idx_blog_article_is_top` (`is_top`)')
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- content: video
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
