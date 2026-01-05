import {defineConfig} from 'vite';
import vue from '@vitejs/plugin-vue';
import {fileURLToPath, URL} from 'node:url';

// Vite 配置：仅负责开发环境的代理配置
// 生产环境的 API 请求由 Nginx 代理处理，文件上传/下载和 WebSocket 的 baseURL 从字典配置中获取
export default defineConfig(({mode}) => {
  const isDev = mode === 'development';

  return {
    // 生产环境部署在 Nginx 的 /admin 路径下
    base: '/admin/',
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: isDev ? 'http://localhost:20000' : 'https://oldbai.top/gateway',
          changeOrigin: true
        }
      }
    }
  };
});

