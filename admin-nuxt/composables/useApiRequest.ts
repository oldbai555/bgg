/**
 * Nuxt 3 API 请求 Composable
 * 符合 Nuxt 3 开发规范，使用 useRuntimeConfig 和 $fetch
 */
export const useApiRequest = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase || 'http://localhost:8888'

  /**
   * 统一处理响应数据
   */
  const handleResponse = (data: any) => {
    // 标准包裹结构：code 为数字（统一错误码）时才按 Envelope 处理
    if (data && typeof data === 'object' && 'code' in data && typeof data.code === 'number') {
      // 支持 code === 0 和 code === 200 作为成功码
      if (data.code === 0 || data.code === 200) {
        return data.data
      }
      const msg = data.msg || '请求失败'
      throw new Error(msg)
    }
    // 非标准结构，直接返回原始 data（兼容字典等特殊接口）
    return data
  }

  /**
   * GET 请求
   */
  const get = async <T = any>(url: string, params?: any): Promise<T> => {
    try {
      const fullUrl = params ? `${baseURL}${url}?${new URLSearchParams(params).toString()}` : `${baseURL}${url}`
      const data = await $fetch(fullUrl, {
        method: 'GET',
        credentials: 'include'
      })
      return handleResponse(data) as T
    } catch (error: any) {
      const msg = error?.data?.msg || error?.message || '请求失败'
      throw new Error(msg)
    }
  }

  /**
   * POST 请求
   */
  const post = async <T = any>(url: string, body?: any): Promise<T> => {
    try {
      const data = await $fetch(`${baseURL}${url}`, {
        method: 'POST',
        body,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      return handleResponse(data) as T
    } catch (error: any) {
      const msg = error?.data?.msg || error?.message || '请求失败'
      throw new Error(msg)
    }
  }

  /**
   * PUT 请求
   */
  const put = async <T = any>(url: string, body?: any): Promise<T> => {
    try {
      const data = await $fetch(`${baseURL}${url}`, {
        method: 'PUT',
        body,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      return handleResponse(data) as T
    } catch (error: any) {
      const msg = error?.data?.msg || error?.message || '请求失败'
      throw new Error(msg)
    }
  }

  /**
   * DELETE 请求
   */
  const del = async <T = any>(url: string, params?: any): Promise<T> => {
    try {
      const fullUrl = params ? `${baseURL}${url}?${new URLSearchParams(params).toString()}` : `${baseURL}${url}`
      const data = await $fetch(fullUrl, {
        method: 'DELETE',
        credentials: 'include'
      })
      return handleResponse(data) as T
    } catch (error: any) {
      const msg = error?.data?.msg || error?.message || '请求失败'
      throw new Error(msg)
    }
  }

  return {
    get,
    post,
    put,
    delete: del
  }
}
