import {
  ping,
  demoList,
  demoCreate,
  demoUpdate,
  demoDelete,
  dailyShortSentenceList,
  dailyShortSentenceCreate,
  dailyShortSentenceUpdate,
  dailyShortSentenceDelete
} from '@/api/generated/admin'

/**
 * Misc 域 API 封装（健康检查 + demo 脚手架示例 + 每日一言）
 * demo/daily_short_sentence 的去留结论见 07-cleanup-and-tooling.md，Week 2 清理时按结论调整本文件。
 */
export const miscApi = {
  ping,
  demoList,
  demoCreate,
  demoUpdate,
  demoDelete,
  dailyShortSentenceList,
  dailyShortSentenceCreate,
  dailyShortSentenceUpdate,
  dailyShortSentenceDelete
}
