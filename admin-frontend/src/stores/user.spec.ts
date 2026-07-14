import {describe, it, expect, beforeEach, vi, afterEach} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {useUserStore} from './user'
import {iamApi} from '@/api/iam'

vi.mock('@/api/iam', () => ({
  iamApi: {
    login: vi.fn(),
    profile: vi.fn(),
    menuMyTree: vi.fn(),
    logout: vi.fn()
  }
}))

vi.mock('@/stores/dict', () => ({
  useDictStore: () => ({
    loadDicts: vi.fn(),
    clearDicts: vi.fn()
  })
}))

vi.mock('@/stores/websocket', () => ({
  useWebSocketStore: () => ({
    connect: vi.fn(),
    disconnect: vi.fn()
  })
}))

describe('useUserStore', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  afterEach(() => {
    localStorage.clear()
  })

  it('login 成功后写入 token 并持久化到 localStorage', async () => {
    vi.mocked(iamApi.login).mockResolvedValue({accessToken: 'at', refreshToken: 'rt'} as never)
    vi.mocked(iamApi.profile).mockResolvedValue({id: 1, username: 'admin', permissions: ['user:list']} as never)
    vi.mocked(iamApi.menuMyTree).mockResolvedValue({list: []} as never)

    const store = useUserStore()
    await store.login({username: 'admin', password: 'x'} as never)

    expect(store.token).toBe('at')
    expect(store.refreshToken).toBe('rt')
    expect(localStorage.getItem('admin_token')).toBe('at')
    expect(store.permissions).toEqual(['user:list'])
  })

  it('cacheValid 在 TTL 内返回 true，过期后返回 false', () => {
    const store = useUserStore()
    store.cacheAt = Date.now()
    expect(store.cacheValid()).toBe(true)

    store.cacheAt = Date.now() - 6 * 60 * 1000 // 超过 5 分钟 TTL
    expect(store.cacheValid()).toBe(false)
  })

  it('cacheAt 为 0（从未缓存）时 cacheValid 返回 false', () => {
    const store = useUserStore()
    store.cacheAt = 0
    expect(store.cacheValid()).toBe(false)
  })

  it('fetchProfile 命中缓存时不重复请求', async () => {
    const store = useUserStore()
    store.cacheAt = Date.now()
    await store.fetchProfile()
    expect(iamApi.profile).not.toHaveBeenCalled()
  })

  it('fetchProfile force=true 时忽略缓存强制请求', async () => {
    vi.mocked(iamApi.profile).mockResolvedValue({id: 1, username: 'admin', permissions: []} as never)
    const store = useUserStore()
    store.cacheAt = Date.now()
    await store.fetchProfile(true)
    expect(iamApi.profile).toHaveBeenCalledOnce()
  })

  it('logout 清空内存状态和 localStorage', async () => {
    const store = useUserStore()
    store.token = 'at'
    store.refreshToken = 'rt'
    store.permissions = ['user:list']
    store.persistCache()
    localStorage.setItem('admin_token', 'at')
    localStorage.setItem('admin_refresh_token', 'rt')

    vi.mocked(iamApi.logout).mockResolvedValue(undefined as never)
    await store.logout()

    expect(store.token).toBe('')
    expect(store.refreshToken).toBe('')
    expect(store.permissions).toEqual([])
    expect(store.menus).toEqual([])
    expect(localStorage.getItem('admin_token')).toBeNull()
    expect(localStorage.getItem('admin_permissions')).toBeNull()
  })

  it('logout 时后端接口报错也不影响本地状态清理', async () => {
    const store = useUserStore()
    store.token = 'at'
    vi.mocked(iamApi.logout).mockRejectedValue(new Error('network error'))

    await store.logout()

    expect(store.token).toBe('')
  })
})
