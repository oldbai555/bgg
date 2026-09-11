import {describe, expect, it} from 'vitest'
import type {MenuItem} from '@/api/generated/admin'
import {findMenuPathChain, resolveActiveMenuPath} from '@/utils/menuActive'

const menus: MenuItem[] = [
  {
    id: 1,
    parentId: 0,
    name: '仪表盘',
    path: '/admin/dashboard',
    component: 'Dashboard',
    icon: '',
    type: 2,
    orderNum: 1,
    visible: 1,
    status: 1
  },
  {
    id: 90,
    parentId: 0,
    name: '博客管理',
    path: '/admin/blog',
    component: '',
    icon: '',
    type: 1,
    orderNum: 2,
    visible: 1,
    status: 1,
    children: [
      {
        id: 93,
        parentId: 90,
        name: '文章管理',
        path: '/admin/blog/article',
        component: 'content/BlogArticleList',
        icon: '',
        type: 2,
        orderNum: 1,
        visible: 1,
        status: 1
      }
    ]
  }
]

describe('resolveActiveMenuPath', () => {
  it('精确匹配菜单 path', () => {
    expect(resolveActiveMenuPath('/admin/dashboard', menus)).toBe('/admin/dashboard')
    expect(resolveActiveMenuPath('/admin/blog/article', menus)).toBe('/admin/blog/article')
  })

  it('子路由对齐到最长前缀菜单', () => {
    expect(resolveActiveMenuPath('/admin/blog/article/edit', menus)).toBe('/admin/blog/article')
    expect(resolveActiveMenuPath('/admin/blog/article/edit/12', menus)).toBe('/admin/blog/article')
  })

  it('无匹配时回退为原 route.path', () => {
    expect(resolveActiveMenuPath('/admin/unknown', menus)).toBe('/admin/unknown')
  })
})

describe('findMenuPathChain', () => {
  it('返回从根到目标的菜单链', () => {
    const chain = findMenuPathChain('/admin/blog/article', menus)
    expect(chain?.map((m) => m.name)).toEqual(['博客管理', '文章管理'])
  })
})
