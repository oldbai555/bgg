import {defineConfig} from 'vitest/config';
import type {ViteDevServer} from 'vite';
import vue from '@vitejs/plugin-vue';
import {fileURLToPath, URL} from 'node:url';

// Vite 配置：仅负责开发环境的代理配置
// 生产环境的 API 请求由 Nginx 代理处理，文件上传/下载和 WebSocket 的 baseURL 从字典配置中获取
export default defineConfig(({mode}) => {
  const isDev = mode === 'development';

  return {
    // 生产环境部署在 Nginx 的 /bgg 路径下（内部再按 /bgg/admin、/bgg/front 分后台/公共两个命名空间）
    // 开发环境也使用 /bgg/ 以保持与生产环境一致
    base: '/bgg/',
    plugins: [
      vue(),
      // 开发环境重定向插件：处理 /bgg 到 /bgg/ 的重定向
      ...(isDev
        ? [
            {
              name: 'redirect-bgg',
              configureServer(server: ViteDevServer) {
                server.middlewares.use((req, res, next) => {
                  // 如果访问 /bgg（不带尾部斜杠），重定向到 /bgg/
                  if (req.url === '/bgg' && !req.url.endsWith('/')) {
                    res.writeHead(301, {Location: '/bgg/'})
                    res.end()
                    return
                  }
                  next()
                })
              }
            }
          ]
        : [])
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    css: {
      preprocessorOptions: {
        scss: {
          // 全局注入间距/圆角/阴影等 SCSS 令牌 + 响应式断点 mixin，业务 .vue 文件无需逐个 @use；
          // 必须用 @use（而非 @import）——Sass 要求 @use 出现在文件其它规则之前，
          // 部分文件（如 layout.scss）已自带 @use './variables.scss' as *，两次 @use 同一模块是幂等的
          additionalData: `@use "@/styles/variables.scss" as *; @use "@/styles/responsive.scss" as *;`
        }
      }
    },
    server: {
      port: 5173,
      // 开发服务器自动打开浏览器时使用正确的路径
      open: '/bgg/front',
      proxy: {
        '/api': {
          target: isDev ? 'http://localhost:20000' : 'https://oldbai.top/gateway',
          changeOrigin: true
        }
      }
    },
    test: {
      // happy-dom 15.x 在本机 Node 26 环境下 window.localStorage 取不到值（已验证是 happy-dom 自身问题，
      // 不是本项目代码问题），按 03 号文档"除非遇到兼容性问题否则优先用 happy-dom"的例外条款切到 jsdom
      environment: 'jsdom',
      globals: true,
      setupFiles: ['./src/test-setup.ts']
    }
  };
});

