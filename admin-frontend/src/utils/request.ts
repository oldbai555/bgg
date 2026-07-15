import axios, {type AxiosResponse, type AxiosError} from 'axios'
import {useUserStore} from '@/stores/user'
import router from '@/router'
import {isEnvelope} from '@/types/envelope'

// API 请求基础地址（仅用于 HTTP API 请求）：
// - 开发环境：通过 Vite dev server 代理到 http://localhost:20000（baseURL 为 /api）
// - 生产环境：浏览器直接请求 /gateway/api/*（由 Nginx 代理到后端）
// 注意：文件上传/下载和 WebSocket 的 baseURL 从字典配置中获取
const baseURL = import.meta.env.PROD ? '/gateway/api' : '/api'

const instance = axios.create({
  baseURL,
  timeout: 15000
})

instance.interceptors.request.use((config) => {
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

// 判断当前路由是否是公共页面（不需要登录）
export const isPublicPath = (): boolean => {
  // 用路由解析后的 path（不含 base 前缀），不要用 window.location.pathname——
  // 生产环境部署在 /bgg/ 下（见 vite.config.ts base），真实 URL 是 /bgg/front/blog/...，
  // 直接判断 pathname.startsWith(...) 永远不成立，会导致公共页游客的 10003 也被强制跳登录页。
  // 公共页统一收在 /front 分支下，和后台 /admin 分支不共享任何路径前缀段。
  const path = router.currentRoute.value.path
  return path.startsWith('/front')
}

// token 失效（10003）时的清理逻辑：先清本地状态，非公共页面才跳登录
export const handleTokenExpired = (): void => {
  const userStore = useUserStore()
  userStore.token = ''
  userStore.refreshToken = ''
  userStore.profile = null
  userStore.permissions = []
  userStore.menus = []
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_refresh_token')
  localStorage.removeItem('admin_permissions')
  localStorage.removeItem('admin_menus')
  localStorage.removeItem('admin_cache_at')

  if (!isPublicPath()) {
    router.push('/admin/login')
  }
}

// 根据后端 Envelope 结构统一处理响应：{ code, msg, data }
export const handleResponse = (resp: {data: unknown}) => {
  const res = resp.data
  // 标准包裹结构：code 为数字（统一错误码）时才按 Envelope 处理
  if (isEnvelope(res)) {
    const code = res.code
    // 支持 code === 0 和 code === 200 作为成功码
    if (code === 0 || code === 200) {
      return res.data
    }

    // 处理 10003 错误码：访问令牌无效或已过期
    if (code === 10003) {
      handleTokenExpired()
    }

    const msg = res.msg || '请求失败'
    return Promise.reject(new Error(msg))
  }
  // 非标准结构，直接返回原始 data（兼容字典等特殊接口）
  return res
}

export const handleResponseError = (error: {
  response?: {data?: {code?: number; msg?: string; message?: string}};
  message?: string;
}) => {
  const data = error?.response?.data
  const code = data?.code

  // 处理 10003 错误码（可能在 error.response.data 中）
  if (code === 10003) {
    handleTokenExpired()
  }

  const msg =
    (data && (data.msg || data.message)) ||
    error.message ||
    '请求失败'
  return Promise.reject(new Error(msg))
}

// handleResponse 故意返回解包后的 res.data 而不是完整 AxiosResponse（项目约定的"响应拦截器直接拆包"），
// 与 axios 声明的拦截器类型不完全一致，这里用 as 桥接；handleResponse/handleResponseError 保持宽松的参数类型是为了单测里能直接构造纯对象调用
instance.interceptors.response.use(
  handleResponse as unknown as (resp: AxiosResponse) => AxiosResponse,
  handleResponseError as unknown as (error: AxiosError) => Promise<never>
)

export default instance

