import {
  // 监控
  monitorStats,
  monitorStatus,
  // 打点
  metricReport,
  metricReportOptions,
  metricStats,
  // 操作日志
  operationLogList,
  operationLogDetail,
  operationLogExport,
  // 登录日志
  loginLogList,
  loginLogDetail,
  loginLogExport,
  loginLogStats,
  // 审计日志
  auditLogList,
  auditLogDetail,
  auditLogExport,
  // 性能日志
  performanceLogList,
  performanceLogExport
} from '@/api/generated/admin'
import type {MetricReportReq, Response} from '@/api/generated/admin'

/**
 * Monitoring 域 API 封装（监控/打点统计/操作日志/登录日志/审计日志/性能日志）
 * 原 src/api/metric.ts 的打点上报/统计能力合并到此处。
 */
export const monitoringApi = {
  // ========== 监控 ==========
  monitorStats,
  monitorStatus,

  // ========== 打点 ==========
  /**
   * 通用打点上报
   * - module: 业务模块标识，如 blog_article_list/blog_article_detail/video_list/video_detail
   * - bizId: 业务ID（文章ID、视频ID等）
   * - event: 事件类型，如 view/play 等
   * - extra: 额外 JSON 字符串
   */
  metricReport: (req: MetricReportReq) => metricReport(req) as Promise<Response>,
  metricReportOptions,
  metricStats,

  // ========== 操作日志 ==========
  operationLogList,
  operationLogDetail,
  operationLogExport,

  // ========== 登录日志 ==========
  loginLogList,
  loginLogDetail,
  loginLogExport,
  loginLogStats,

  // ========== 审计日志 ==========
  auditLogList,
  auditLogDetail,
  auditLogExport,

  // ========== 性能日志 ==========
  performanceLogList,
  performanceLogExport
}
