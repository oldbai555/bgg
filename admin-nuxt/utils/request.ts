/**
 * 注意：在 Nuxt 3 中，推荐使用 $fetch 或 useFetch 进行 API 请求
 * 此文件保留用于兼容性，但建议使用 composables/useApiRequest.ts
 * 
 * 如果必须使用 axios，应该在 composable 中创建实例
 */
import axios from 'axios'

// 注意：useRuntimeConfig 只能在 setup 函数或 composable 中调用
// 这里提供一个工厂函数，在 composable 中使用
export const createAxiosInstance = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase || 'http://localhost:8888'

  const instance = axios.create({
    baseURL,
    timeout: 15000
  })

  // 请求拦截器
  instance.interceptors.request.use((config) => {
    // 如果需要 token，可以在这里添加
    // const token = useCookie('token')
    // if (token.value) {
    //   config.headers = config.headers || {}
    //   config.headers.Authorization = `Bearer ${token.value}`
    // }
    return config
  })

  // 响应拦截器
  // 根据后端 Envelope 结构统一处理响应：{ code, msg, data }
  instance.interceptors.response.use(
    (resp) => {
      const res = resp.data
      // 标准包裹结构：code 为数字（统一错误码）时才按 Envelope 处理
      if (res && typeof res === 'object' && 'code' in res && typeof (res as any).code === 'number') {
        // 支持 code === 0 和 code === 200 作为成功码
        if ((res as any).code === 0 || (res as any).code === 200) {
          return (res as any).data
        }
        const msg = (res as any).msg || '请求失败'
        return Promise.reject(new Error(msg))
      }
      // 非标准结构，直接返回原始 data（兼容字典等特殊接口）
      return res
    },
    (error) => {
      const data = error?.response?.data
      const msg =
        (data && (data.msg || data.message)) ||
        error.message ||
        '请求失败'
      return Promise.reject(new Error(msg))
    }
  )

  return instance
}

// 默认导出（用于向后兼容，但不推荐）
// 注意：这会在非 setup 上下文中失败
let defaultInstance: ReturnType<typeof createAxiosInstance> | null = null

export default new Proxy({} as ReturnType<typeof createAxiosInstance>, {
  get(_target, prop) {
    if (!defaultInstance) {
      // 延迟初始化，在首次使用时创建
      defaultInstance = createAxiosInstance()
    }
    return (defaultInstance as any)[prop]
  }
})
