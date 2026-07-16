/**
 * 应用配置管理 Composable
 * 统一从字典获取配置，避免硬编码
 */
import {ref, computed} from 'vue'
import {systemApi} from '@/api/system'

// 配置缓存
const configCache = ref<Record<string, string>>({})
const configLoading = ref<Record<string, boolean>>({})

/**
 * 从字典获取配置值
 * @param code 字典类型编码
 * @param defaultValue 默认值（如果字典中没有配置）
 * @returns 配置值
 */
export async function getDictConfig(code: string, defaultValue: string = ''): Promise<string> {
  // 如果正在加载，等待
  if (configLoading.value[code]) {
    return new Promise((resolve) => {
      const checkInterval = setInterval(() => {
        if (!configLoading.value[code]) {
          clearInterval(checkInterval)
          resolve(configCache.value[code] || defaultValue)
        }
      }, 100)
    })
  }

  // 如果缓存中有，直接返回
  if (configCache.value[code]) {
    return configCache.value[code]
  }

  // 从字典加载配置
  configLoading.value[code] = true
  try {
    const resp = await systemApi.dictGet({code})
    if (resp && resp.items && resp.items.length > 0) {
      const value = resp.items[0].value
      configCache.value[code] = value
      return value
    }
    return defaultValue
  } catch (err) {
    console.warn(`获取字典配置失败: ${code}`, err)
    return defaultValue
  } finally {
    configLoading.value[code] = false
  }
}

/**
 * 获取存储 baseURL（文件上传/下载）
 */
export async function getStorageBaseURL(): Promise<string> {
  return await getDictConfig('storage_base_url', '')
}

/**
 * 获取 WebSocket baseURL
 */
export async function getWebSocketBaseURL(): Promise<string> {
  return await getDictConfig('websocket_base_url', '')
}

/**
 * 清除配置缓存（当字典配置更新后调用）
 */
export function clearConfigCache(code?: string) {
  if (code) {
    delete configCache.value[code]
  } else {
    configCache.value = {}
  }
}

/**
 * 使用配置的 Composable
 */
export function useAppConfig() {
  // 存储 baseURL（响应式）
  const storageBaseURL = ref<string>('')
  // WebSocket baseURL（响应式）
  const websocketBaseURL = ref<string>('')

  // 初始化配置
  const initConfig = async () => {
    storageBaseURL.value = await getStorageBaseURL()
    websocketBaseURL.value = await getWebSocketBaseURL()
  }

  // 刷新配置
  const refreshConfig = async () => {
    clearConfigCache()
    await initConfig()
  }

  return {
    storageBaseURL: computed(() => storageBaseURL.value),
    websocketBaseURL: computed(() => websocketBaseURL.value),
    initConfig,
    refreshConfig
  }
}

