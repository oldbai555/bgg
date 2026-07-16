import {describe, it, expect, beforeEach} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {useAppStore} from './app'

describe('useAppStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('setTheme 更新 state 并持久化到 localStorage 与 DOM 属性', () => {
    const store = useAppStore()
    store.setTheme('dark')

    expect(store.theme).toBe('dark')
    expect(localStorage.getItem('admin_theme')).toBe('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })

  it('setLang 更新 state 并持久化', () => {
    const store = useAppStore()
    store.setLang('en')

    expect(store.lang).toBe('en')
    expect(localStorage.getItem('admin_lang')).toBe('en')
  })

  it('toggleSidebar 翻转折叠状态并持久化', () => {
    const store = useAppStore()
    expect(store.sidebarCollapsed).toBe(false)

    store.toggleSidebar()
    expect(store.sidebarCollapsed).toBe(true)
    expect(localStorage.getItem('admin_sidebar_collapsed')).toBe('true')

    store.toggleSidebar()
    expect(store.sidebarCollapsed).toBe(false)
    expect(localStorage.getItem('admin_sidebar_collapsed')).toBe('false')
  })

  it('setSidebarCollapsed 直接设置指定状态', () => {
    const store = useAppStore()
    store.setSidebarCollapsed(true)
    expect(store.sidebarCollapsed).toBe(true)
    expect(localStorage.getItem('admin_sidebar_collapsed')).toBe('true')
  })

  it('init() 把当前 theme 写入 DOM data-theme 属性', () => {
    localStorage.setItem('admin_theme', 'dark')
    const store = useAppStore()
    store.init()

    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })
})
