import {taskList, taskCancel, taskDetail, taskRecent} from '@/api/generated/admin'

/**
 * Task 域 API 封装（异步任务）
 */
export const taskApi = {
  taskList,
  taskCancel,
  taskDetail,
  taskRecent
}
