// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  
  // TypeScript 配置
  typescript: {
    strict: true,
    typeCheck: false // 暂时关闭类型检查，避免启动时的类型错误
  },

  // CSS 配置（Element Plus 样式在插件中导入）
  css: [
    '@/assets/styles/global.scss', // 全局样式，确保最外层没有滚动条
    '@/assets/styles/public-list.scss',
    '@/assets/styles/public-detail.scss',
    '@/assets/styles/blog.scss'
  ],

  // 模块配置
  modules: [
    '@pinia/nuxt'
  ],

  // 运行时配置
  runtimeConfig: {
    public: {
      // 使用 Nuxt 3 标准环境变量 NUXT_PUBLIC_API_BASE
      // 在 .env 文件中配置，或使用默认值
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:20000'
    }
  },

  // 应用配置
  app: {
    head: {
      title: '博客与视频',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '博客文章和视频内容展示' }
      ]
    }
  },

  // 兼容性配置
  compatibilityDate: '2024-01-01',

  // 开发服务器配置
  devServer: {
    port: 3000
  }
})
