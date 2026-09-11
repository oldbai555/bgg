import type {MenuItem} from '@/api/generated/admin'

/**
 * 在菜单树中找与当前路由最匹配的菜单 path。
 * 只考虑 type=2（实际 el-menu-item）；优先精确匹配，否则取最长前缀。
 * 编辑页 /admin/blog/article/edit/:id 应对齐列表菜单 /admin/blog/article。
 */
export function resolveActiveMenuPath(routePath: string, menus: MenuItem[]): string {
  if (!routePath) {
    return ''
  }

  let best = ''
  const walk = (items: MenuItem[]) => {
    for (const item of items) {
      const p = item.path
      // type=2 才是侧栏可高亮的菜单项；目录 type=1 的 path 通常不是 el-menu-item index
      if (item.type === 2 && p) {
        if (routePath === p || routePath.startsWith(p + '/')) {
          if (p.length > best.length) {
            best = p
          }
        }
      }
      if (item.children?.length) {
        walk(item.children)
      }
    }
  }
  walk(menus)
  return best || routePath
}

/**
 * 从菜单根走到目标 path 的祖先链（含目标自身）。
 * targetPath 应为 resolveActiveMenuPath 的结果或菜单上真实存在的 path。
 */
export function findMenuPathChain(targetPath: string, menus: MenuItem[]): MenuItem[] | null {
  const walk = (items: MenuItem[], chain: MenuItem[]): MenuItem[] | null => {
    for (const item of items) {
      const next = [...chain, item]
      if (item.path === targetPath) {
        return next
      }
      if (item.children?.length) {
        const found = walk(item.children, next)
        if (found) {
          return found
        }
      }
    }
    return null
  }
  return walk(menus, [])
}
