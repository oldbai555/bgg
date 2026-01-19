import axios from 'axios';
import {useUserStore} from '@/stores/user';
import router from '@/router';

// API 请求基础地址（仅用于 HTTP API 请求）：
// - 开发环境：通过 Vite dev server 代理到 http://localhost:20000（baseURL 为 /api）
// - 生产环境：浏览器直接请求 /gateway/api/*（由 Nginx 代理到后端）
// 注意：文件上传/下载和 WebSocket 的 baseURL 从字典配置中获取
const baseURL = import.meta.env.PROD ? '/gateway/api' : '/api';

const instance = axios.create({
  baseURL,
  timeout: 15000
});

instance.interceptors.request.use((config) => {
  const userStore = useUserStore();
  if (userStore.token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${userStore.token}`;
  }
  return config;
});

// 判断当前路由是否是公共页面（不需要登录）
const isPublicPath = (): boolean => {
  if (typeof window === 'undefined') {
    return false;
  }
  const path = window.location.pathname;
  return path.startsWith('/blog') || path.startsWith('/videos');
};

// 根据后端 Envelope 结构统一处理响应：{ code, msg, data }
instance.interceptors.response.use(
  (resp) => {
    const res = resp.data;
    // 标准包裹结构：code 为数字（统一错误码）时才按 Envelope 处理
    if (res && typeof res === 'object' && 'code' in res && typeof (res as any).code === 'number') {
      const code = (res as any).code;
      // 支持 code === 0 和 code === 200 作为成功码
      if (code === 0 || code === 200) {
        return (res as any).data;
      }
      
      // 处理 10003 错误码：访问令牌无效或已过期
      if (code === 10003) {
        const userStore = useUserStore();
        // 先清除本地状态，避免发送不必要的请求
        userStore.token = '';
        userStore.refreshToken = '';
        userStore.profile = null;
        userStore.permissions = [];
        userStore.menus = [];
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_refresh_token');
        localStorage.removeItem('admin_permissions');
        localStorage.removeItem('admin_menus');
        localStorage.removeItem('admin_cache_at');
        
        // 如果不是公共页面，才跳转到登录页
        if (!isPublicPath()) {
          router.push('/login');
        }
      }
      
      const msg = (res as any).msg || '请求失败';
      return Promise.reject(new Error(msg));
    }
    // 非标准结构，直接返回原始 data（兼容字典等特殊接口）
    return res;
  },
  (error) => {
    const data = error?.response?.data;
    const code = data?.code;
    
    // 处理 10003 错误码（可能在 error.response.data 中）
    if (code === 10003) {
      const userStore = useUserStore();
      // 先清除本地状态，避免发送不必要的请求
      userStore.token = '';
      userStore.refreshToken = '';
      userStore.profile = null;
      userStore.permissions = [];
      userStore.menus = [];
      localStorage.removeItem('admin_token');
      localStorage.removeItem('admin_refresh_token');
      localStorage.removeItem('admin_permissions');
      localStorage.removeItem('admin_menus');
      localStorage.removeItem('admin_cache_at');
      
      // 如果不是公共页面，才跳转到登录页
      if (!isPublicPath()) {
        router.push('/login');
      }
    }
    
    const msg =
      (data && (data.msg || data.message)) ||
      error.message ||
      '请求失败';
    return Promise.reject(new Error(msg));
  }
);

export default instance;

