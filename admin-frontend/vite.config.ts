import {defineConfig} from 'vite';
import vue from '@vitejs/plugin-vue';
import {fileURLToPath, URL} from 'node:url';

// 使用函数形式，根据运行模式切换配置
export default defineConfig(({mode}) => {
  // dev 环境：本地后端
  // prod 环境：线上网关
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
          target: isDev ? 'http://localhost:8888' : 'https://oldbai.top/gateway',
          changeOrigin: true
        }
      }
    }
  };
});

