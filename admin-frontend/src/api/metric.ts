import {metricReport, metricStats} from '@/api/generated/admin';
import type {MetricReportReq, MetricStatsReq, MetricStatsResp, Response} from '@/api/generated/admin';

/**
 * 通用打点上报
 * - module: 业务模块标识，如 blog_article_list/blog_article_detail/video_list/video_detail
 * - bizId: 业务ID（文章ID、视频ID等）
 * - event: 事件类型，如 view/play 等
 * - extra: 额外 JSON 字符串
 */
export const metricApi = {
  report: (req: MetricReportReq) => metricReport(req) as Promise<Response>,
  stats: (req: MetricStatsReq) => metricStats(req) as Promise<MetricStatsResp>
};

