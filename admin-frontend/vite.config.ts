import {defineConfig} from 'vite';
import vue from '@vitejs/plugin-vue';
import {fileURLToPath, URL} from 'node:url';

// Vite 配置：仅负责开发环境的代理配置
// 生产环境的 API 请求由 Nginx 代理处理，文件上传/下载和 WebSocket 的 baseURL 从字典配置中获取
export default defineConfig(({mode}) => {
  const isDev = mode === 'development';

  return {
    // 生产环境部署在 Nginx 的 /admin 路径下
    // 开发环境也使用 /admin/ 以保持与生产环境一致
    base: '/admin/',
    plugins: [
      vue(),
      // 开发环境重定向插件：处理 /admin 到 /admin/ 的重定向
      ...(isDev
        ? [
            {
              name: 'redirect-admin',
              configureServer(server) {
                server.middlewares.use((req, res, next) => {
                  // 如果访问 /admin（不带尾部斜杠），重定向到 /admin/
                  if (req.url === '/admin' && !req.url.endsWith('/')) {
                    res.writeHead(301, {Location: '/admin/'})
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
    server: {
      port: 5173,
      // 开发服务器自动打开浏览器时使用正确的路径
      open: '/admin/',
      proxy: {
        '/api': {
          target: isDev ? 'http://localhost:20000' : 'https://oldbai.top/gateway',
          changeOrigin: true
        }
      }
    }
  };
});

