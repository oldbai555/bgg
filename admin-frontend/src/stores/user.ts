import {defineStore} from 'pinia'
import {iamApi} from '@/api/iam'
import type {LoginReq, ProfileResp, MenuItem} from '@/api/generated/admin'

interface UserState {
  token: string;
  refreshToken: string;
  profile: ProfileResp | null;
  permissions: string[];
  menus: MenuItem[];
  cacheAt: number;
}

const tokenKey = 'admin_token'
const refreshKey = 'admin_refresh_token'
const permKey = 'admin_permissions'
const menuKey = 'admin_menus'
const cacheAtKey = 'admin_cache_at'
const CACHE_TTL = 5 * 60 * 1000

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem(tokenKey) || '',
    refreshToken: localStorage.getItem(refreshKey) || '',
    profile: null,
    permissions: JSON.parse(localStorage.getItem(permKey) || '[]'),
    menus: JSON.parse(localStorage.getItem(menuKey) || '[]'),
    cacheAt: Number(localStorage.getItem(cacheAtKey) || 0)
  }),
  actions: {
    async login(payload: LoginReq) {
      const data = await iamApi.login(payload)
      await this.afterLoginSuccess(data.accessToken, data.refreshToken)
    },
    async loginByFeishu(code: string, state: string) {
      const data = await iamApi.loginFeishu({code, state})
      await this.afterLoginSuccess(data.accessToken, data.refreshToken)
    },
    // afterLoginSuccess 拿到 token 后的公共收尾：存 token、拉 profile/menus/字典、连 WebSocket。
    // 密码登录和飞书登录殊途同归，都走这一段，避免重复维护两套初始化逻辑。
    async afterLoginSuccess(accessToken: string, refreshToken: string) {
      this.token = accessToken
      this.refreshToken = refreshToken
      localStorage.setItem(tokenKey, this.token)
      localStorage.setItem(refreshKey, this.refreshToken)
      await this.fetchProfile(true)
      await this.fetchMenus(true)

      // 登录后加载字典数据
      const {useDictStore} = await import('./dict')
      const dictStore = useDictStore()
      await dictStore.loadDicts()

      // 登录后自动连接 WebSocket（如果有权限）
      const {useWebSocketStore} = await import('./websocket')
      const wsStore = useWebSocketStore()
      wsStore.connect()
    },
    async fetchProfile(force = false) {
      if (!force && this.cacheValid()) {
        return
      }
      const profileData = await iamApi.profile()
      this.profile = profileData
      this.permissions = profileData.permissions || []
      this.persistCache()
    },
    async fetchMenus(force = false) {
      if (!force && this.cacheValid() && this.menus.length > 0) {
        return
      }
      // 使用 my-tree 接口，根据用户权限过滤菜单
      const resp = await iamApi.menuMyTree()
      this.menus = resp.list || []
      this.persistCache()
    },
    async logout() {
      try {
        await iamApi.logout({
          accessToken: this.token,
          refreshToken: this.refreshToken
        })
      } catch {
        // ignore
      }

      // 退出登录时断开 WebSocket
      const {useWebSocketStore} = await import('./websocket')
      const wsStore = useWebSocketStore()
      wsStore.disconnect()

      this.token = ''
      this.refreshToken = ''
      this.profile = null
      this.permissions = []
      this.menus = []
      this.cacheAt = 0
      localStorage.removeItem(tokenKey)
      localStorage.removeItem(refreshKey)
      localStorage.removeItem(permKey)
      localStorage.removeItem(menuKey)
      localStorage.removeItem(cacheAtKey)

      // 退出登录时清除字典数据
      const {useDictStore} = await import('./dict')
      const dictStore = useDictStore()
      dictStore.clearDicts()
    },
    cacheValid() {
      if (!this.cacheAt) {
        return false
      }
      return Date.now() - this.cacheAt < CACHE_TTL
    },
    persistCache() {
      this.cacheAt = Date.now()
      localStorage.setItem(permKey, JSON.stringify(this.permissions || []))
      localStorage.setItem(menuKey, JSON.stringify(this.menus || []))
      localStorage.setItem(cacheAtKey, String(this.cacheAt))
    }
  }
})

