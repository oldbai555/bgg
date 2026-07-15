import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import permissionDirective from './directives/permission'
import {useAppStore} from './stores/app'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './styles/theme.scss'
import './styles/layout.scss'

// 全局错误处理：忽略浏览器扩展相关的错误和 Vite HMR 错误
window.addEventListener('error', (event) => {
  const errorMessage = event.message || event.filename || event.error?.message || ''
  // 忽略浏览器扩展相关的错误
  if (
    errorMessage.includes('message channel closed') ||
    errorMessage.includes('asynchronous response') ||
    errorMessage.includes('Extension context invalidated') ||
    errorMessage.includes('runtime.lastError') ||
    errorMessage.includes('message port closed') ||
    errorMessage.includes('The message port closed before a response was received')
  ) {
    event.preventDefault()
    return false
  }
  // 忽略 Vite HMR 相关的错误（开发环境）
  if (
    import.meta.env.DEV &&
    (errorMessage.includes('/src/router/index.ts') ||
     errorMessage.includes('ERR_ABORTED') ||
     errorMessage.includes('500'))
  ) {
    // 开发环境下的 HMR 错误可以忽略
    event.preventDefault()
    return false
  }
})

// 处理未捕获的 Promise 错误
window.addEventListener('unhandledrejection', (event) => {
  // 忽略浏览器扩展相关的错误
  const errorMessage = event.reason?.message || event.reason?.toString() || ''
  if (
    errorMessage.includes('message channel closed') ||
    errorMessage.includes('asynchronous response') ||
    errorMessage.includes('Extension context invalidated') ||
    errorMessage.includes('runtime.lastError') ||
    errorMessage.includes('message port closed') ||
    errorMessage.includes('The message port closed before a response was received')
  ) {
    event.preventDefault()
    return false
  }
  // 忽略 Vite HMR 相关的错误（开发环境）
  if (
    import.meta.env.DEV &&
    (errorMessage.includes('/src/router/index.ts') ||
     errorMessage.includes('ERR_ABORTED') ||
     errorMessage.includes('500'))
  ) {
    event.preventDefault()
    return false
  }
})

// 拦截 console.error，过滤浏览器扩展相关的警告
const originalConsoleError = console.error
console.error = (...args: unknown[]) => {
  const message = args.join(' ')
  // 忽略浏览器扩展相关的错误
  if (
    message.includes('runtime.lastError') ||
    message.includes('The message port closed before a response was received') ||
    message.includes('message channel closed') ||
    message.includes('Extension context invalidated')
  ) {
    // 静默忽略，不输出到控制台
    return
  }
  // 其他错误正常输出
  originalConsoleError.apply(console, args)
}

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(i18n)
app.use(ElementPlus, {
  locale: zhCn
})
app.directive('permission', permissionDirective)

const appStore = useAppStore(pinia)
appStore.init()

if (appStore.lang) {
  i18n.global.locale.value = appStore.lang
}

app.mount('#app')

