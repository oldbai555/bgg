import {describe, it, expect, beforeEach, vi} from 'vitest'
import {systemApi} from '@/api/system'
import {
  getDictConfig,
  getStorageBaseURL,
  getWebSocketBaseURL,
  clearConfigCache,
  useAppConfig
} from './useAppConfig'

vi.mock('@/api/system', () => ({
  systemApi: {
    dictGet: vi.fn()
  }
}))

describe('useAppConfig', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    clearConfigCache()
  })

  it('getDictConfig 命中字典时返回第一个 item 的 value 并写入缓存', async () => {
    vi.mocked(systemApi.dictGet).mockResolvedValue({items: [{value: 'https://cdn.example.com'}]} as never)

    const value = await getDictConfig('storage_base_url', 'fallback')

    expect(value).toBe('https://cdn.example.com')
    expect(systemApi.dictGet).toHaveBeenCalledWith({code: 'storage_base_url'})

    // 命中缓存后不再重复请求
    const cached = await getDictConfig('storage_base_url', 'fallback')
    expect(cached).toBe('https://cdn.example.com')
    expect(systemApi.dictGet).toHaveBeenCalledTimes(1)
  })

  it('getDictConfig 字典为空或请求失败时回退默认值，不向上抛错', async () => {
    vi.mocked(systemApi.dictGet).mockResolvedValueOnce({items: []} as never)
    expect(await getDictConfig('empty_code', 'default1')).toBe('default1')

    vi.mocked(systemApi.dictGet).mockRejectedValueOnce(new Error('network error'))
    expect(await getDictConfig('error_code', 'default2')).toBe('default2')
  })

  it('getStorageBaseURL / getWebSocketBaseURL 透传各自的字典 code', async () => {
    vi.mocked(systemApi.dictGet).mockImplementation((async (req: {code: string}) => ({
      items: [{value: `value-of-${req.code}`}]
    })) as never)

    expect(await getStorageBaseURL()).toBe('value-of-storage_base_url')
    expect(await getWebSocketBaseURL()).toBe('value-of-websocket_base_url')
  })

  it('clearConfigCache 按 code 清除或清空全部，之后重新请求字典', async () => {
    vi.mocked(systemApi.dictGet).mockResolvedValue({items: [{value: 'v1'}]} as never)
    await getDictConfig('a', '')
    await getDictConfig('b', '')
    expect(systemApi.dictGet).toHaveBeenCalledTimes(2)

    clearConfigCache('a')
    await getDictConfig('a', '')
    await getDictConfig('b', '')
    // a 缓存被清了要重新请求，b 命中缓存不请求
    expect(systemApi.dictGet).toHaveBeenCalledTimes(3)
  })

  it('useAppConfig().initConfig 并行加载 storage/websocket baseURL 并通过 computed 暴露', async () => {
    vi.mocked(systemApi.dictGet).mockImplementation((async (req: {code: string}) => ({
      items: [{value: `${req.code}-value`}]
    })) as never)

    const {storageBaseURL, websocketBaseURL, initConfig} = useAppConfig()
    expect(storageBaseURL.value).toBe('')

    await initConfig()

    expect(storageBaseURL.value).toBe('storage_base_url-value')
    expect(websocketBaseURL.value).toBe('websocket_base_url-value')
  })

  it('useAppConfig().refreshConfig 清缓存后重新拉取', async () => {
    vi.mocked(systemApi.dictGet).mockResolvedValue({items: [{value: 'first'}]} as never)
    const {storageBaseURL, initConfig, refreshConfig} = useAppConfig()
    await initConfig()
    expect(storageBaseURL.value).toBe('first')

    vi.mocked(systemApi.dictGet).mockResolvedValue({items: [{value: 'second'}]} as never)
    await refreshConfig()
    expect(storageBaseURL.value).toBe('second')
  })
})
