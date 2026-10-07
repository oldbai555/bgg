import {describe, expect, it} from 'vitest'
import {
  addDays,
  buildLinePath,
  canCheckinOn,
  diffDays,
  formatRate,
  fromDayDraft,
  localDateStr,
  suggestTrainingStatus,
  toDayDraft,
  valueRange,
  weekdayOf
} from './fitness'
import {FitnessTrainingStatus} from '@/constants/fitness'

describe('fitness utils', () => {
  it('addDays 跨月跨年', () => {
    expect(addDays('2026-10-07', 1)).toBe('2026-10-08')
    expect(addDays('2026-10-01', -1)).toBe('2026-09-30')
    expect(addDays('2026-12-31', 1)).toBe('2027-01-01')
    expect(addDays('2028-02-28', 1)).toBe('2028-02-29')
  })

  it('diffDays / weekdayOf', () => {
    expect(diffDays('2026-10-07', '2026-09-30')).toBe(7)
    expect(weekdayOf('2026-10-05')).toBe(1)
    expect(weekdayOf('2026-10-11')).toBe(7)
  })

  it('canCheckinOn 只允许今天和最近 7 天', () => {
    expect(canCheckinOn('2026-10-07', '2026-10-07', 7)).toBe(true)
    expect(canCheckinOn('2026-09-30', '2026-10-07', 7)).toBe(true)
    expect(canCheckinOn('2026-09-29', '2026-10-07', 7)).toBe(false)
    expect(canCheckinOn('2026-10-08', '2026-10-07', 7)).toBe(false)
  })

  it('训练状态推荐与完成率格式化', () => {
    expect(suggestTrainingStatus(0, 5)).toBe(FitnessTrainingStatus.None)
    expect(suggestTrainingStatus(2, 5)).toBe(FitnessTrainingStatus.Partial)
    expect(suggestTrainingStatus(5, 5)).toBe(FitnessTrainingStatus.Done)
    expect(suggestTrainingStatus(1, 0)).toBe(FitnessTrainingStatus.None)
    expect(formatRate(0.756)).toBe('76%')
  })

  it('localDateStr 按本机日历日补零', () => {
    expect(localDateStr(new Date(2026, 0, 5))).toBe('2026-01-05')
    expect(localDateStr(new Date(2026, 11, 31))).toBe('2026-12-31')
  })

  it('buildLinePath 缺值处断开', () => {
    expect(buildLinePath([], 100, 50, 0, 10)).toBe('')
    expect(buildLinePath([10, 0, 0, 5], 30, 10, 0, 10)).toBe('M0,0 M30,5')
    expect(buildLinePath([0, 5, 10], 20, 10, 0, 10)).toBe('M10,5 L20,0')
  })

  it('周计划编辑态互转：空内容给默认值，提交时去掉空动作/空餐次', () => {
    const draft = toDayDraft(undefined)
    expect(draft.trainingType).toBe(3)
    expect(draft.exercises).toEqual([])

    draft.trainingType = 1
    draft.exercises.push({name: ' 深蹲 ', sets: 3, reps: '8-12'}, {name: '  '})
    draft.meals.push({slot: 'lunch', food: '鸡腿饭'}, {slot: ''})
    const content = fromDayDraft(draft)
    expect(content.exercises).toEqual([{name: '深蹲', sets: 3, reps: '8-12'}])
    expect(content.meals).toEqual([{slot: 'lunch', food: '鸡腿饭'}])

    const back = toDayDraft(content)
    back.exercises[0].name = '硬拉'
    expect(content.exercises?.[0].name).toBe('深蹲')
  })

  it('valueRange 忽略缺值并留余量', () => {
    expect(valueRange([0, 0])).toEqual([0, 1])
    expect(valueRange([70.2, 0, 71.6])).toEqual([69, 73])
  })
})
