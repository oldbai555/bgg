import {publicVideoList, publicVideoDetail} from '@/api/generated/admin'
import type {
  PublicVideoListReq,
  PublicVideoListResp,
  PublicVideoDetailReq,
  PublicVideoDetailResp
} from '@/api/generated/admin'

/**
 * 视频相关 API 封装
 */
export const videoApi = {
  // 公开视频列表
  publicList: (req: PublicVideoListReq) =>
    publicVideoList(req) as Promise<PublicVideoListResp>,

  // 公开视频详情
  publicDetail: (req: PublicVideoDetailReq) =>
    publicVideoDetail(req) as Promise<PublicVideoDetailResp>
}
