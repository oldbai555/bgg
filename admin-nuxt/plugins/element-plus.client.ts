/**
 * Element Plus 插件（仅客户端）
 * 在 Nuxt 3 中注册 Element Plus 组件和样式
 */
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

export default defineNuxtPlugin((nuxtApp) => {
  // 注册 Element Plus（带中文语言包）
  nuxtApp.vueApp.use(ElementPlus, {
    locale: zhCn
  })

  // 注册所有图标
  for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    nuxtApp.vueApp.component(key, component)
  }
})
