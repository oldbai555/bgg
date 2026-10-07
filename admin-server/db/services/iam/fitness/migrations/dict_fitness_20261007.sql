-- ============================================
-- 字典SQL增量脚本
-- 模块：身材管理（fitness）
-- 创建时间：2026-10-07
-- 说明：数值型枚举与 services/iam/internal/consts/fitness.go 一一对应，value 从 1 开始，
--       0 表示「未填/不筛选」（如训练状态 0=当天还没打训练卡）。
--       另给 metric_module 字典追加手机端 5 个页面（含登录页）的埋点模块（字符串 value，同 dict_metric_module_20260714.sql）。
--       admin_dict_item 没有业务唯一键，字典项一律 INSERT ... SELECT ... WHERE NOT EXISTS，可重复执行。
-- ============================================

INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('身材管理训练类型', 'fitness_training_type', '模板每日训练类型', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理训练打卡状态', 'fitness_training_status', '训练打卡状态，0 表示未打卡', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理饮食打卡状态', 'fitness_meal_status', '每餐打卡状态，0 表示未打卡', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理训练时间偏好', 'fitness_training_preference', '早练/晚练', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理模板切换来源', 'fitness_switch_source', '模板切换记录的来源', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理模板建议状态', 'fitness_suggestion_status', '每周模板建议的处理状态', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('身材管理餐次', 'fitness_meal_slot', '一天五个餐次，value 是周模板/打卡 meals JSON 里的 slot 编码', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
SELECT dt.`id`, s.`label`, s.`value`, s.`sort`, 1, s.`remark`, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
FROM (
  SELECT 'fitness_training_type' AS `code`, '力量' AS `label`, '1' AS `value`, 1 AS `sort`, 'FitnessTrainingStrength' AS `remark`
  UNION ALL SELECT 'fitness_training_type', '有氧', '2', 2, 'FitnessTrainingCardio'
  UNION ALL SELECT 'fitness_training_type', '休息', '3', 3, 'FitnessTrainingRest'
  UNION ALL SELECT 'fitness_training_status', '完成', '1', 1, 'FitnessTrainingStatusDone'
  UNION ALL SELECT 'fitness_training_status', '部分完成', '2', 2, 'FitnessTrainingStatusPartial'
  UNION ALL SELECT 'fitness_training_status', '跳过', '3', 3, 'FitnessTrainingStatusSkipped'
  UNION ALL SELECT 'fitness_meal_status', '按计划吃了', '1', 1, 'FitnessMealStatusOnPlan'
  UNION ALL SELECT 'fitness_meal_status', '吃了别的', '2', 2, 'FitnessMealStatusOther'
  UNION ALL SELECT 'fitness_meal_status', '没吃', '3', 3, 'FitnessMealStatusSkipped'
  UNION ALL SELECT 'fitness_training_preference', '早练', '1', 1, 'FitnessPreferenceMorning'
  UNION ALL SELECT 'fitness_training_preference', '晚练', '2', 2, 'FitnessPreferenceEvening'
  UNION ALL SELECT 'fitness_switch_source', '手动切换', '1', 1, 'FitnessSwitchSourceManual'
  UNION ALL SELECT 'fitness_switch_source', '采纳建议', '2', 2, 'FitnessSwitchSourceSuggestion'
  UNION ALL SELECT 'fitness_switch_source', '管理员指定', '3', 3, 'FitnessSwitchSourceAdmin'
  UNION ALL SELECT 'fitness_switch_source', '首次登录默认', '4', 4, 'FitnessSwitchSourceDefault'
  UNION ALL SELECT 'fitness_suggestion_status', '待处理', '1', 1, 'FitnessSuggestionPending'
  UNION ALL SELECT 'fitness_suggestion_status', '已采纳', '2', 2, 'FitnessSuggestionAccepted'
  UNION ALL SELECT 'fitness_suggestion_status', '已拒绝', '3', 3, 'FitnessSuggestionRejected'
  UNION ALL SELECT 'fitness_suggestion_status', '已过期', '4', 4, 'FitnessSuggestionExpired'
  UNION ALL SELECT 'fitness_meal_slot', '早餐', 'breakfast', 1, 'FitnessMealBreakfast'
  UNION ALL SELECT 'fitness_meal_slot', '午餐', 'lunch', 2, 'FitnessMealLunch'
  UNION ALL SELECT 'fitness_meal_slot', '加餐', 'snack', 3, 'FitnessMealSnack'
  UNION ALL SELECT 'fitness_meal_slot', '晚餐', 'dinner', 4, 'FitnessMealDinner'
  UNION ALL SELECT 'fitness_meal_slot', '到家后', 'after_home', 5, 'FitnessMealAfterHome'
  UNION ALL SELECT 'metric_module', '身材管理-今日计划', 'fitness_today', 11, 'MetricModuleFitnessToday'
  UNION ALL SELECT 'metric_module', '身材管理-打卡', 'fitness_checkin', 12, 'MetricModuleFitnessCheckin'
  UNION ALL SELECT 'metric_module', '身材管理-身体数据', 'fitness_body', 13, 'MetricModuleFitnessBody'
  UNION ALL SELECT 'metric_module', '身材管理-我的', 'fitness_me', 14, 'MetricModuleFitnessMe'
  UNION ALL SELECT 'metric_module', '身材管理-登录', 'fitness_login', 15, 'MetricModuleFitnessLogin'
) s
JOIN `admin_dict_type` dt ON dt.`code` = s.`code` AND dt.`deleted_at` = 0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_dict_item` di
  WHERE di.`type_id` = dt.`id` AND di.`value` = s.`value` AND di.`deleted_at` = 0
);
