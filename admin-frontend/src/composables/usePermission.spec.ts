import {describe, it, expect, beforeEach} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {usePermission} from './usePermission'
import {useUserStore} from '@/stores/user'

describe('usePermission', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('未传 code 时始终放行', () => {
    const {hasPermission} = usePermission()
    expect(hasPermission()).toBe(true)
    expect(hasPermission('')).toBe(true)
  })

  it('权限列表为空时拒绝任何具体 code', () => {
    const {hasPermission} = usePermission()
    expect(hasPermission('user:list')).toBe(false)
  })

  it('拥有通配符 * 时放行任意 code', () => {
    const userStore = useUserStore()
    userStore.permissions = ['*']
    const {hasPermission} = usePermission()
    expect(hasPermission('user:list')).toBe(true)
    expect(hasPermission('anything:else')).toBe(true)
  })

  it('命中具体 code 时放行，未命中时拒绝', () => {
    const userStore = useUserStore()
    userStore.permissions = ['user:list', 'role:list']
    const {hasPermission} = usePermission()
    expect(hasPermission('user:list')).toBe(true)
    expect(hasPermission('user:delete')).toBe(false)
  })
})
