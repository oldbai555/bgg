import request from '@/utils/request'
import type {DictGetReq, DictGetResp} from '@/api/generated/admin'

/**
 * 公共字典查询（未登录可用，仅白名单 code，例如 video_proxy_url）
 * 走 /api/v1/public/dict，而不是受鉴权的 /api/v1/dict
 * 注意：request 的 baseURL 已包含 /api，所以这里只需要 /v1/public/dict
 */
export function publicDictGet(params: DictGetReq) {
  return request.get<DictGetResp>('/v1/public/dict', {params})
}

