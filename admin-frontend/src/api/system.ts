import {
  // 文件
  fileList,
  fileCreate,
  fileUpdate,
  fileDelete,
  fileDownload,
  fileUpload,
  // 配置
  configList,
  configCreate,
  configUpdate,
  configDelete,
  configGet,
  // 字典类型/字典项（鉴权后台管理接口；未登录可用的字典查询走 @/api/public.ts）
  dictTypeList,
  dictTypeCreate,
  dictTypeUpdate,
  dictTypeDelete,
  dictItemList,
  dictItemCreate,
  dictItemUpdate,
  dictItemDelete,
  dictGet,
  dictBatchGet,
  // 公告
  noticeList,
  noticeCreate,
  noticeUpdate,
  noticeDelete,
  // 通知
  notificationList,
  notificationDelete,
  notificationRead,
  notificationClearRead,
  notificationReadAll
} from '@/api/generated/admin'

/**
 * System 域 API 封装（配置/字典/文件/公告/通知）
 */
export const systemApi = {
  // ========== 文件 ==========
  fileList,
  fileCreate,
  fileUpdate,
  fileDelete,
  fileDownload,
  fileUpload,

  // ========== 配置 ==========
  configList,
  configCreate,
  configUpdate,
  configDelete,
  configGet,

  // ========== 字典类型 ==========
  dictTypeList,
  dictTypeCreate,
  dictTypeUpdate,
  dictTypeDelete,

  // ========== 字典项 ==========
  dictItemList,
  dictItemCreate,
  dictItemUpdate,
  dictItemDelete,

  // ========== 字典查询（鉴权态） ==========
  dictGet,
  dictBatchGet,

  // ========== 公告 ==========
  noticeList,
  noticeCreate,
  noticeUpdate,
  noticeDelete,

  // ========== 通知 ==========
  notificationList,
  notificationDelete,
  notificationRead,
  notificationClearRead,
  notificationReadAll
}
