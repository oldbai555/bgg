import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import {ElMessage} from 'element-plus'
import {useUserStore} from '@/stores/user'
import {usePermission} from '@/composables/usePermission'
import type {MenuItem} from '@/api/generated/admin'

const viewModules = import.meta.glob('../views/**/*.vue')
const knownViewKeys = new Set(
  Object.keys(viewModules).map((key) => key.replace(/^\.\.\/views\//, '').replace(/\.vue$/, ''))
)

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/front'
  },
  // 公共页面（不需要登录）
  {
    path: '/front',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/front/blog',
    name: 'BlogList',
    component: () => import('@/views/public/BlogList.vue')
  },
  {
    path: '/front/blog/:id',
    name: 'BlogDetail',
    component: () => import('@/views/public/BlogDetail.vue')
  },
  {
    path: '/front/videos',
    name: 'VideoList',
    component: () => import('@/views/public/VideoList.vue')
  },
  {
    path: '/front/videos/:id',
    name: 'VideoDetail',
    component: () => import('@/views/public/VideoDetail.vue')
  },
  {
    path: '/admin/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/layout',
    name: 'Root',
    component: () => import('@/layouts/DefaultLayout.vue'),
    children: [
      {
        path: '/admin/dashboard',
        name: 'Dashboard',
        meta: {keepAlive: true},
        component: () => import('@/views/Dashboard.vue')
      },
      // 博客文章编辑页（不挂菜单，通过文章列表跳转进入）
      {
        path: '/admin/blog/article/edit',
        name: 'BlogArticleCreate',
        meta: {permission: 'blog_article:create', keepAlive: false},
        component: () => import('@/views/content/BlogArticleEdit.vue')
      },
      {
        path: '/admin/blog/article/edit/:id',
        name: 'BlogArticleEdit',
        meta: {permission: 'blog_article:update', keepAlive: false},
        component: () => import('@/views/content/BlogArticleEdit.vue')
      },
      {
        path: '/admin/system/role',
        name: 'RoleList',
        meta: {permission: 'role:list', keepAlive: true},
        component: () => import('@/views/iam/RoleList.vue')
      },
      {
        path: '/admin/system/permission',
        name: 'PermissionList',
        meta: {permission: 'permission:list', keepAlive: true},
        component: () => import('@/views/iam/PermissionList.vue')
      },
      {
        path: '/admin/system/department',
        name: 'DepartmentList',
        meta: {permission: 'department:tree', keepAlive: true},
        component: () => import('@/views/iam/DepartmentList.vue')
      },
      {
        path: '/admin/system/api',
        name: 'ApiList',
        meta: {permission: 'api:list', keepAlive: true},
        component: () => import('@/views/iam/ApiList.vue')
      },
      {
        path: '/admin/system/profile',
        name: 'Profile',
        meta: {keepAlive: true},
        component: () => import('@/views/iam/Profile.vue')
      },
      {
        path: '/admin/system/task',
        name: 'TaskListPage',
        meta: {permission: 'task:list', keepAlive: true},
        component: () => import('@/views/task/TaskList.vue')
      },
      {
        path: '/admin/system/metric-stats',
        name: 'MetricStats',
        meta: {permission: 'metric:stats', keepAlive: true},
        component: () => import('@/views/monitoring/MetricStats.vue')
      },
      {
        path: '/admin/403',
        name: 'NoAccess',
        component: () => import('@/views/error/NoAccess.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory('/bgg/'),
  routes
})

// 添加全局错误处理
router.onError((error) => {
  console.error('[Router] 路由错误:', error)
})

let initialized = false
const dynamicAdded = new Set<string>()

function resolveComponent(component?: string, path?: string) {
  const candidates: string[] = []
  if (component) {
    candidates.push(component.replace(/^\//, ''))
  }
  if (path) {
    candidates.push(path.replace(/^\//, ''))
  }
  for (const [key, loader] of Object.entries(viewModules)) {
    const clean = key.replace(/^..\/views\//, '').replace(/\.vue$/, '')
    if (candidates.includes(clean)) {
      return loader
    }
  }
  // 如果找不到组件，输出错误信息
  if (component || path) {
    console.error(`[Router] 无法解析组件: component="${component}", path="${path}"`)
  }
  return undefined
}

export function generateUniqueRouteName(rawPath: string, usedNames: Set<string>): string {
  const base = rawPath.replace(/^\//, '').replace(/\//g, '_') || 'root'
  let uniqueName = base
  let counter = 1
  while (usedNames.has(uniqueName)) {
    uniqueName = `${base}_${counter}`
    counter++
  }
  usedNames.add(uniqueName)
  return uniqueName
}

function buildRoutesFromMenus(menus: MenuItem[]): RouteRecordRaw[] {
  const res: RouteRecordRaw[] = []
  const usedNames = new Set<string>()

  const walk = (items: MenuItem[]) => {
    items.forEach((m) => {
      // 处理菜单（type=2）：有实际页面组件
      if (m.type === 2 && m.path) {
        const comp = resolveComponent(m.component, m.path)
        if (comp) {
          const uniqueName = generateUniqueRouteName(m.path, usedNames)

          res.push({
            path: m.path,
            name: uniqueName,
            meta: {
              permission: m.permissionCode || undefined, // 如果没有权限码，设为 undefined
              keepAlive: true
            },
            component: comp
          })
        } else {
          console.error(`[Router] 路由注册失败: path="${m.path}", component="${m.component}"`)
        }
      }
      // 处理目录（type=1）：如果有子菜单，重定向到第一个子菜单
      if (m.type === 1 && m.path && m.children && m.children.length > 0) {
        // 找到第一个有效的子菜单
        const firstChild = m.children.find((child) => child.type === 2 && child.path)
        if (firstChild && firstChild.path) {
          const uniqueName = generateUniqueRouteName(m.path, usedNames)

          res.push({
            path: m.path,
            name: uniqueName,
            redirect: firstChild.path,
            meta: {
              permission: m.permissionCode || undefined // 如果没有权限码，设为 undefined
            }
          })
        }
      }
      if (m.children?.length) {
        walk(m.children)
      }
    })
  }
  walk(menus)
  return res
}

let menusValidated = false

// dev 环境启动期校验：把菜单 component 与 views/ 目录的失配从"点击后静默 404"提前到"登录后立即可见"
function validateMenuComponents(menus: MenuItem[]) {
  if (!import.meta.env.DEV || menusValidated) {
    return
  }
  menusValidated = true

  const mismatched: string[] = []
  const walk = (items: MenuItem[]) => {
    items.forEach((m) => {
      if (m.type === 2 && m.component) {
        const clean = m.component.replace(/^\//, '')
        if (!knownViewKeys.has(clean)) {
          mismatched.push(`「${m.name}」component="${m.component}"`)
        }
      }
      if (m.children?.length) {
        walk(m.children)
      }
    })
  }
  walk(menus)

  if (mismatched.length > 0) {
    console.warn(`[Router] 以下菜单的 component 在 src/views 下找不到对应文件，页面将无法访问：\n${mismatched.join('\n')}`)
    ElMessage.warning(`发现 ${mismatched.length} 个菜单的 component 路径失配，详见控制台`)
  }
}

router.beforeEach(async (to, _from, next) => {
  try {
    const userStore = useUserStore()
    const {hasPermission} = usePermission()

    // 公共页面不需要登录：路由方案里公共页统一收在 /front 分支下，
    // 和后台 /admin 分支不共享任何路径前缀段，两者永不冲突
    const isPublicPath = to.path.startsWith('/front')

    if (!isPublicPath && to.path !== '/admin/login' && !userStore.token) {
      next('/admin/login')
      return
    }

    if (userStore.token) {
      // 如果未初始化或菜单数据为空，重新获取
      if (!initialized || !userStore.menus || userStore.menus.length === 0) {
        initialized = true
        try {
          await userStore.fetchProfile()
        } catch (err) {
          console.error('[Router] 获取用户信息失败:', err)
        }
        try {
          await userStore.fetchMenus()
          validateMenuComponents(userStore.menus)
        } catch (err) {
          console.error('[Router] 获取菜单失败:', err)
        }
      }

      // 添加动态路由
      if (userStore.menus && userStore.menus.length > 0) {
        const dynRoutes = buildRoutesFromMenus(userStore.menus)
        let addedNew = false
        dynRoutes.forEach((r) => {
          if (!dynamicAdded.has(r.path as string)) {
            try {
              router.addRoute('Root', r)
              dynamicAdded.add(r.path as string)
              addedNew = true
            } catch (err) {
              console.error(`[Router] 添加路由失败: ${r.path}`, err)
            }
          }
        })

        // 首次加载（尤其是硬刷新）时，本次导航可能是在动态路由注册前就已经解析完成的，
        // 无论当时解析到的是 NotFound 还是某个碰巧匹配上但并非目标页面的静态路由，
        // 只要这一轮确实新注册了路由，就用 to.fullPath 重新解析一次，解析结果和当前不一致就纠正过去，
        // 避免只处理 NotFound 这一种情况而漏掉"匹配到错误页面"的情况
        if (addedNew) {
          const resolved = router.resolve(to.fullPath)
          if (resolved.name && resolved.name !== to.name) {
            next({...resolved, replace: true})
            return
          }
        }
      }
    }

    // 权限检查：只有当 meta.permission 存在且不为空时才检查权限
    const needPerm = to.meta?.permission as string | undefined
    if (needPerm && needPerm.trim() !== '' && !hasPermission(needPerm)) {
      next('/admin/403')
      return
    }

    next()
  } catch (err) {
    console.error('[Router] 路由守卫错误:', err)
    // 如果路由守卫出错，仍然允许导航（避免阻塞）
    next()
  }
})

export default router
