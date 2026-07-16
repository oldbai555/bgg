import {describe, it, expect, beforeEach, vi} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'

const {currentRoute} = vi.hoisted(() => ({
  currentRoute: {value: {path: '/admin/dashboard'}}
}))

vi.mock('@/router', () => ({
  default: {push: vi.fn(), currentRoute}
}))

import router from '@/router'
import {handleResponse, handleResponseError, isPublicPath} from './request'
import {useUserStore} from '@/stores/user'

const setPath = (path: string) => {
  currentRoute.value.path = path
}

describe('isPublicPath', () => {
  it('/front 开头的路径判定为公共页面', () => {
    setPath('/front/blog/1')
    expect(isPublicPath()).toBe(true)
    setPath('/front/videos')
    expect(isPublicPath()).toBe(true)
  })

  it('后台管理路径判定为非公共页面', () => {
    setPath('/admin/dashboard')
    expect(isPublicPath()).toBe(false)
  })

  it('用路由解析后的 path 判断，不受生产环境 /bgg base 前缀影响', () => {
    // router.currentRoute.value.path 本身就是去掉 base 前缀之后的路径（Vue Router 行为），
    // 这里模拟的正是这种"真实 URL 是 /bgg/front/blog，但 route.path 是 /front/blog"的场景
    setPath('/front/blog/1')
    expect(isPublicPath()).toBe(true)
  })
})

describe('handleResponse', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    setPath('/admin/dashboard')
  })

  it('code 为 0 时解包并返回 data', () => {
    const result = handleResponse({data: {code: 0, msg: 'ok', data: {id: 1}}})
    expect(result).toEqual({id: 1})
  })

  it('code 为 200 时同样解包并返回 data（兼容两种成功码）', () => {
    const result = handleResponse({data: {code: 200, msg: 'ok', data: {id: 2}}})
    expect(result).toEqual({id: 2})
  })

  it('非 Envelope 结构时原样返回（兼容字典等特殊接口）', () => {
    const raw = {list: [1, 2, 3]}
    const result = handleResponse({data: raw})
    expect(result).toBe(raw)
  })

  it('业务错误码（非 0/200）时 reject 并携带后端 msg', async () => {
    await expect(
      Promise.resolve(handleResponse({data: {code: 10001, msg: '参数错误', data: null}}))
    ).rejects.toThrow('参数错误')
  })

  it('code 10003（token 过期）时清空登录态并跳转登录页', async () => {
    const userStore = useUserStore()
    userStore.token = 'at'
    userStore.permissions = ['user:list']
    localStorage.setItem('admin_token', 'at')

    await expect(
      Promise.resolve(handleResponse({data: {code: 10003, msg: '登录已过期', data: null}}))
    ).rejects.toThrow('登录已过期')

    expect(userStore.token).toBe('')
    expect(userStore.permissions).toEqual([])
    expect(localStorage.getItem('admin_token')).toBeNull()
    expect(router.push).toHaveBeenCalledWith('/admin/login')
  })

  it('公共页面上 10003 不跳转登录页（避免打断游客浏览）', async () => {
    setPath('/front/blog/1')
    await expect(
      Promise.resolve(handleResponse({data: {code: 10003, msg: '登录已过期', data: null}}))
    ).rejects.toThrow()
    expect(router.push).not.toHaveBeenCalled()
  })
})

describe('handleResponseError', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    setPath('/admin/dashboard')
  })

  it('优先使用 response.data.msg 作为错误信息', async () => {
    await expect(
      handleResponseError({response: {data: {msg: '服务器错误'}}})
    ).rejects.toThrow('服务器错误')
  })

  it('response.data.msg 缺失时回退到 error.message', async () => {
    await expect(
      handleResponseError({message: 'Network Error'})
    ).rejects.toThrow('Network Error')
  })

  it('什么信息都没有时回退到默认文案', async () => {
    await expect(handleResponseError({})).rejects.toThrow('请求失败')
  })

  it('error.response.data.code 为 10003 时同样清空登录态', async () => {
    const userStore = useUserStore()
    userStore.token = 'at'

    await expect(
      handleResponseError({response: {data: {code: 10003, msg: '登录已过期'}}})
    ).rejects.toThrow()

    expect(userStore.token).toBe('')
    expect(router.push).toHaveBeenCalledWith('/admin/login')
  })
})
