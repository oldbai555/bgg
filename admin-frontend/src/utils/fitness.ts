/**
 * 身材管理前端纯函数。日期一律是后端给的 Asia/Shanghai 日历日字符串（YYYY-MM-DD），
 * 这里只做字符串层面的加减/比较，用 UTC 计算避免受浏览器本地时区影响。
 */
import type {FitnessDayContent, FitnessExercise, FitnessMeal} from '@/api/generated/admin'
import {FitnessTrainingStatus, FitnessTrainingType} from '@/constants/fitness'

const DAY_MS = 24 * 60 * 60 * 1000
const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

const toUTC = (date: string): number => {
  const [y, m, d] = date.split('-').map(Number)
  return Date.UTC(y, m - 1, d)
}

const fromUTC = (ms: number): string => {
  const dt = new Date(ms)
  const m = String(dt.getUTCMonth() + 1).padStart(2, '0')
  const d = String(dt.getUTCDate()).padStart(2, '0')
  return `${dt.getUTCFullYear()}-${m}-${d}`
}

/** 日期选择器给的本地 Date → YYYY-MM-DD（选择器按本机日历日工作） */
export const localDateStr = (d: Date): string => {
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${d.getFullYear()}-${m}-${day}`
}

export const addDays = (date: string, n: number): string => fromUTC(toUTC(date) + n * DAY_MS)

/** a - b 相差的天数 */
export const diffDays = (a: string, b: string): number => Math.round((toUTC(a) - toUTC(b)) / DAY_MS)

/** 1=周一 ... 7=周日 */
export const weekdayOf = (date: string): number => {
  const w = new Date(toUTC(date)).getUTCDay()
  return w === 0 ? 7 : w
}

export const weekdayLabel = (weekday: number): string => WEEKDAY_LABELS[weekday - 1] || ''

/** 补卡窗口：今天及往前 backfillDays 天 */
export const canCheckinOn = (date: string, today: string, backfillDays: number): boolean => {
  const d = diffDays(today, date)
  return d >= 0 && d <= backfillDays
}

export const formatRate = (rate: number): string => `${Math.round((rate || 0) * 100)}%`

/** 根据勾选情况推荐训练状态：全勾=完成，勾了部分=部分完成，一个没勾=未打卡（由用户自己选跳过） */
export const suggestTrainingStatus = (done: number, total: number): number => {
  if (total <= 0 || done <= 0) {
    return FitnessTrainingStatus.None
  }
  return done >= total ? FitnessTrainingStatus.Done : FitnessTrainingStatus.Partial
}

/** 折线图坐标：把一组可能缺值（<=0）的数值映射到 [0,width]x[0,height]，返回 SVG path */
export const buildLinePath = (values: number[], width: number, height: number, min: number, max: number): string => {
  if (values.length === 0) {
    return ''
  }
  const span = max - min || 1
  const step = values.length > 1 ? width / (values.length - 1) : 0
  let path = ''
  let pen = false
  values.forEach((v, i) => {
    if (!v || v <= 0) {
      pen = false
      return
    }
    const x = +(i * step).toFixed(1)
    const y = +(height - ((v - min) / span) * height).toFixed(1)
    path += `${pen ? 'L' : 'M'}${x},${y} `
    pen = true
  })
  return path.trim()
}

/** 后台周计划编辑态：数组/字符串字段都给默认值，便于表单直接 v-model 和 push/splice */
export interface FitnessDayDraft {
  trainingType: number;
  trainingTime: string;
  stepGoal: number;
  warmup: string;
  exercises: FitnessExercise[];
  stretch: string;
  trainingNote: string;
  meals: FitnessMeal[];
  tip: string;
}

export const toDayDraft = (content?: FitnessDayContent | null): FitnessDayDraft => ({
  trainingType: content?.trainingType || FitnessTrainingType.Rest,
  trainingTime: content?.trainingTime || '',
  stepGoal: content?.stepGoal || 0,
  warmup: content?.warmup || '',
  exercises: (content?.exercises || []).map((e) => ({...e})),
  stretch: content?.stretch || '',
  trainingNote: content?.trainingNote || '',
  meals: (content?.meals || []).map((m) => ({...m, alternatives: [...(m.alternatives || [])]})),
  tip: content?.tip || ''
})

/** 去掉没填名字的动作、没选餐次的饮食，其余字段原样提交（后端再做长度/合法性校验） */
export const fromDayDraft = (draft: FitnessDayDraft): FitnessDayContent => ({
  ...draft,
  exercises: draft.exercises.filter((e) => e.name.trim()).map((e) => ({...e, name: e.name.trim()})),
  meals: draft.meals.filter((m) => m.slot)
})

/** 取一组数值里大于 0 的最小/最大值，并上下各留 pad 的余量 */
export const valueRange = (values: number[], pad = 0.5): [number, number] => {
  const valid = values.filter((v) => v > 0)
  if (valid.length === 0) {
    return [0, 1]
  }
  return [Math.floor(Math.min(...valid) - pad), Math.ceil(Math.max(...valid) + pad)]
}
