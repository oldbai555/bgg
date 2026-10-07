-- iam/fitness 建表（身材管理：计划模板 / 打卡 / 身体数据 / 模板切换与建议）
-- fitness 域并入 iam-rpc（admin 库），理由见 docs/changelog/2026-10-07.md：
-- 使用者身份就是 admin_user + admin_user_third_party 的飞书绑定，同库才能直接 JOIN 用户昵称，
-- 且没有独立扩缩容/独立信任边界的理由（15-service-boundaries.md 第 6 节「不拆第 6 个数据服务」）。
--
-- 日期字段统一 VARCHAR(10) 存 Asia/Shanghai 的 'YYYY-MM-DD'：字典序即时间序，按天唯一键直观，
-- 也避免 DATE 列在不同时区 MySQL 会话下被隐式换算。
-- 唯一键都带上 deleted_at，软删后可以重新插入同一业务键。

CREATE TABLE IF NOT EXISTS `fitness_template` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `code` VARCHAR(32) NOT NULL COMMENT '模板编码，切换建议规则按编码识别（starter/standard/advanced/evening/busy/plateau/maintain）',
  `name` VARCHAR(64) NOT NULL COMMENT '模板名称',
  `level` INT NOT NULL DEFAULT 0 COMMENT '升降级序列：1入门 2标准 3进阶，0表示不在升降级序列里',
  `stage` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '建议使用阶段，如 第1-4周',
  `summary` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '模板说明',
  `daily_kcal` INT NOT NULL DEFAULT 0 COMMENT '每日建议摄入热量(千卡)',
  `daily_protein` INT NOT NULL DEFAULT 0 COMMENT '每日建议蛋白质(克)',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_template_code` (`code`, `deleted_at`),
  KEY `idx_fitness_template_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-计划模板';

CREATE TABLE IF NOT EXISTS `fitness_template_day` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_id` BIGINT UNSIGNED NOT NULL COMMENT '所属模板 fitness_template.id',
  `weekday` TINYINT NOT NULL DEFAULT 0 COMMENT '周几：1周一~7周日；0表示按具体日期覆盖',
  `plan_date` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '按日期覆盖的日期 YYYY-MM-DD；周模板为空串',
  `training_type` TINYINT NOT NULL DEFAULT 3 COMMENT '训练类型：1力量 2有氧 3休息',
  `training_time` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '训练时间，如 08:00-09:00',
  `step_goal` INT NOT NULL DEFAULT 0 COMMENT '步数目标',
  `warmup` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '热身',
  `exercises` TEXT NOT NULL COMMENT '动作清单 JSON：[{name,sets,reps,duration,note}]',
  `stretch` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '拉伸',
  `training_note` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '训练要点',
  `meals` TEXT NOT NULL COMMENT '三餐 JSON：[{slot,time,place,food,portion,kcal,protein,howToOrder,alternatives}]',
  `tip` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '当天专属提示，和通用提示（fitness_tip）一起展示，可为空',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_template_day` (`template_id`, `weekday`, `plan_date`, `deleted_at`),
  KEY `idx_fitness_template_day_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-模板每日内容（周模板 + 按日期覆盖）';

CREATE TABLE IF NOT EXISTS `fitness_tip` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `code` VARCHAR(32) NOT NULL COMMENT '提示编码，如 clothes',
  `title` VARCHAR(64) NOT NULL COMMENT '标题',
  `content` TEXT NOT NULL COMMENT '内容，一行一条',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_tip_code` (`code`, `deleted_at`),
  KEY `idx_fitness_tip_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-通用提示';

CREATE TABLE IF NOT EXISTS `fitness_profile` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '使用者 admin_user.id',
  `height_cm` DECIMAL(5,1) NOT NULL DEFAULT 0 COMMENT '身高(cm)',
  `target_weight_kg` DECIMAL(5,1) NOT NULL DEFAULT 0 COMMENT '目标体重(kg)',
  `training_preference` TINYINT NOT NULL DEFAULT 1 COMMENT '训练时间偏好：1早练 2晚练',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)，即首次使用时间',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_profile_user` (`user_id`, `deleted_at`),
  KEY `idx_fitness_profile_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-使用者档案';

CREATE TABLE IF NOT EXISTS `fitness_template_switch` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '使用者 admin_user.id',
  `from_template_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '原模板，0表示首次分配',
  `to_template_id` BIGINT UNSIGNED NOT NULL COMMENT '新模板',
  `effective_date` VARCHAR(10) NOT NULL COMMENT '生效日期 YYYY-MM-DD（含当天）',
  `source` TINYINT NOT NULL DEFAULT 1 COMMENT '来源：1手动切换 2采纳建议 3管理员指定 4首次登录默认',
  `reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '切换原因',
  `operator_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作人 admin_user.id',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_fitness_template_switch_user_date` (`user_id`, `effective_date`),
  KEY `idx_fitness_template_switch_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-模板切换记录（只追加，某天用哪个模板 = 生效日期<=当天的最后一条）';

CREATE TABLE IF NOT EXISTS `fitness_checkin` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '使用者 admin_user.id',
  `checkin_date` VARCHAR(10) NOT NULL COMMENT '打卡日期 YYYY-MM-DD',
  `template_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '打卡时当天计划所属模板（快照）',
  `training_type` TINYINT NOT NULL DEFAULT 3 COMMENT '打卡时当天训练类型（快照）：1力量 2有氧 3休息',
  `training_status` TINYINT NOT NULL DEFAULT 0 COMMENT '训练状态：0未打卡 1完成 2部分完成 3跳过',
  `exercise_done` TEXT NOT NULL COMMENT '已勾选动作下标 JSON 数组',
  `exercise_total` INT NOT NULL DEFAULT 0 COMMENT '打卡时动作总数（快照）',
  `meals` TEXT NOT NULL COMMENT '三餐打卡 JSON：[{slot,status,note}]，status 1按计划 2吃了别的 3没吃',
  `steps` INT NOT NULL DEFAULT 0 COMMENT '步数',
  `step_goal` INT NOT NULL DEFAULT 0 COMMENT '打卡时步数目标（快照）',
  `water_cups` INT NOT NULL DEFAULT 0 COMMENT '喝水杯数（250ml/杯）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_checkin_user_date` (`user_id`, `checkin_date`, `deleted_at`),
  KEY `idx_fitness_checkin_date` (`checkin_date`),
  KEY `idx_fitness_checkin_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-每日打卡';

CREATE TABLE IF NOT EXISTS `fitness_body_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '使用者 admin_user.id',
  `record_date` VARCHAR(10) NOT NULL COMMENT '记录日期 YYYY-MM-DD',
  `weight_kg` DECIMAL(5,1) NOT NULL DEFAULT 0 COMMENT '体重(kg)，0表示未填',
  `waist_cm` DECIMAL(5,1) NOT NULL DEFAULT 0 COMMENT '腰围(cm)，0表示未填',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_body_record_user_date` (`user_id`, `record_date`, `deleted_at`),
  KEY `idx_fitness_body_record_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-身体数据（仅本人与后台管理员可见）';

CREATE TABLE IF NOT EXISTS `fitness_suggestion` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '使用者 admin_user.id',
  `week_start` VARCHAR(10) NOT NULL COMMENT '建议所属周的周一 YYYY-MM-DD',
  `rule_code` VARCHAR(32) NOT NULL COMMENT '命中的规则编码',
  `from_template_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '生成建议时的模板',
  `to_template_id` BIGINT UNSIGNED NOT NULL COMMENT '建议切换到的模板',
  `reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '建议理由（给用户看）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1待处理 2已采纳 3已拒绝 4已过期',
  `decided_at` BIGINT NOT NULL DEFAULT 0 COMMENT '用户处理时间(秒级时间戳)',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(秒级时间戳)',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(秒级时间戳)',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间(秒级时间戳,0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fitness_suggestion_user_week` (`user_id`, `week_start`, `deleted_at`),
  KEY `idx_fitness_suggestion_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身材管理-每周模板建议（只建议不强切）';
