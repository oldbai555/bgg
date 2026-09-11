import {RouteLocationNormalized} from 'vue-router'
import type {MenuItem} from '@/api/generated/admin'
import type {BreadcrumbItem} from '@/components/layout/Breadcrumb.vue'
import {findMenuPathChain, resolveActiveMenuPath} from '@/utils/menuActive'

/**
 * 从菜单数据生成面包屑
 */
export function generateBreadcrumb(
  route: RouteLocationNormalized,
  menus: MenuItem[]
): BreadcrumbItem[] {
  const breadcrumbs: BreadcrumbItem[] = []

  // 首页
  breadcrumbs.push({
    title: '首页',
    path: '/admin/dashboard'
  })

  const activePath = resolveActiveMenuPath(route.path, menus)
  const menuPath = findMenuPathChain(activePath, menus)
  if (menuPath) {
    menuPath.forEach((menu) => {
      if (menu.path && menu.type === 2) {
        breadcrumbs.push({
          title: menu.name,
          path: menu.path
        })
      }
    })
  } else {
    // 如果没有找到对应的菜单，使用路由路径生成
    const pathSegments = route.path.split('/').filter(Boolean)
    pathSegments.forEach((segment, index) => {
      const path = '/' + pathSegments.slice(0, index + 1).join('/')
      breadcrumbs.push({
        title: segment.charAt(0).toUpperCase() + segment.slice(1),
        path: index === pathSegments.length - 1 ? undefined : path
      })
    })
  }

  return breadcrumbs
}
