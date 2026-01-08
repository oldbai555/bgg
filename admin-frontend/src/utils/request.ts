import axios from 'axios';
import {useUserStore} from '@/stores/user';

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

// 根据后端 Envelope 结构统一处理响应：{ code, msg, data }
instance.interceptors.response.use(
  (resp) => {
    const res = resp.data;
    // 标准包裹结构：code 为数字（统一错误码）时才按 Envelope 处理
    if (res && typeof res === 'object' && 'code' in res && typeof (res as any).code === 'number') {
      // 支持 code === 0 和 code === 200 作为成功码
      if ((res as any).code === 0 || (res as any).code === 200) {
        return (res as any).data;
      }
      const msg = (res as any).msg || '请求失败';
      return Promise.reject(new Error(msg));
    }
    // 非标准结构，直接返回原始 data（兼容字典等特殊接口）
    return res;
  },
  (error) => {
    const data = error?.response?.data;
    const msg =
      (data && (data.msg || data.message)) ||
      error.message ||
      '请求失败';
    return Promise.reject(new Error(msg));
  }
);

export default instance;

