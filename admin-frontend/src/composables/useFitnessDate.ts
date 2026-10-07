import {ref, watch} from 'vue'

const STORAGE_KEY = 'fitness_selected_date'

interface StoredDate {
  date: string;
  today: string;
}

const readStored = (): StoredDate | null => {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    return raw ? (JSON.parse(raw) as StoredDate) : null
  } catch {
    return null
  }
}

/**
 * 今日计划页和打卡页共享「当前查看的日期」：从计划页点「去打卡」或切 tab 回来时保持同一天。
 * 空字符串表示「今天」，今天由后端按 Asia/Shanghai 决定，前端不自己算；
 * 存储时连同当时的 today 一起存，跨天后（飞书 webview 可能常驻）自动回到今天。
 */
export function useFitnessDate() {
  const stored = readStored()
  const date = ref(stored?.date || '')
  const today = ref(stored?.today || '')

  watch([date, today], ([d, t]) => {
    if (d && t) {
      sessionStorage.setItem(STORAGE_KEY, JSON.stringify({date: d, today: t}))
    }
  })

  /** 用接口返回的 today 校准；返回 true 表示存储的日期已跨天作废，调用方应按今天重新加载 */
  const syncToday = (serverToday: string): boolean => {
    const stale = !!today.value && today.value !== serverToday
    today.value = serverToday
    if (stale) {
      date.value = serverToday
    }
    return stale
  }

  return {date, today, syncToday}
}
