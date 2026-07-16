import {
  ping,
  dailyShortSentenceList,
  dailyShortSentenceCreate,
  dailyShortSentenceUpdate,
  dailyShortSentenceDelete
} from '@/api/generated/admin'

/**
 * Misc 域 API 封装（健康检查 + 每日一言）
 * demo 相关接口去留结论见 07-cleanup-and-tooling.md：DemoList.vue 是开发流程示例脚手架，Week 2 已删除，本文件同步移除对应导出。
 */
export const miscApi = {
  ping,
  dailyShortSentenceList,
  dailyShortSentenceCreate,
  dailyShortSentenceUpdate,
  dailyShortSentenceDelete
}
