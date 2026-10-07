/**
 * 身材管理手机端常量。数值枚举与后端 services/iam/internal/consts/fitness.go、
 * 字典 dict_fitness_20261007.sql 一一对应，页面上展示的文字一律走字典，这里只放判断逻辑要用的值。
 */
export const FITNESS_HOME_PATH = '/front/fitness'
export const FITNESS_LOGIN_PATH = '/front/fitness/login'
export const FITNESS_CHECKIN_PATH = '/front/fitness/checkin'
export const FITNESS_BODY_PATH = '/front/fitness/body'
export const FITNESS_ME_PATH = '/front/fitness/me'
export const FITNESS_ADMIN_STATS_PATH = '/admin/fitness/stats'

/** 允许补打卡的天数（不含今天），与后端 FitnessBackfillDays 一致 */
export const FITNESS_BACKFILL_DAYS = 7

export const FitnessTrainingType = {
  Strength: 1,
  Cardio: 2,
  Rest: 3
} as const

export const FitnessTrainingStatus = {
  None: 0,
  Done: 1,
  Partial: 2,
  Skipped: 3
} as const

export const FitnessMealStatus = {
  None: 0,
  OnPlan: 1,
  Other: 2,
  Skipped: 3
} as const

export const FitnessTrainingPreference = {
  Morning: 1,
  Evening: 2
} as const

export const FitnessSuggestionStatus = {
  Pending: 1
} as const

/** 模板/通用提示的启用状态（实体自身状态列，不走字典） */
export const FitnessStatus = {
  Enabled: 1,
  Disabled: 2
} as const

/** 餐次固定顺序，展示文字走字典 fitness_meal_slot */
export const FITNESS_MEAL_SLOTS = ['breakfast', 'lunch', 'snack', 'dinner', 'after_home'] as const

export const FITNESS_DICT_CODES = [
  'fitness_training_type',
  'fitness_training_status',
  'fitness_meal_status',
  'fitness_meal_slot',
  'fitness_training_preference',
  'fitness_switch_source',
  'fitness_suggestion_status'
] as const

export const isFitnessPath = (path: string): boolean => path.startsWith(FITNESS_HOME_PATH)
