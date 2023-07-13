import type {AxiosInstance, AxiosResponse, InternalAxiosRequestConfig} from 'axios';
import axios from 'axios'
import {message} from 'ant-design-vue';// 可以替换自己想要的库,这边自动生成就用这个了

const axiosInstance: AxiosInstance = axios.create({
    baseURL: 'http://xxx.xxx.xxx/gateway',
    timeout: 5000,
});

// 设置post请求头
axiosInstance.defaults.headers.post['Content-Type'] = 'application/json';


// 添加请求拦截器
axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // 登录流程控制中，根据本地是否存在token判断用户的登录情况
        // 但是即使token存在，也有可能token是过期的，所以在每次的请求头中携带token
        // 后台根据携带的token判断用户的登录情况，并返回给我们对应的状态码
        const token = localStorage.getItem('ACCESS_TOKEN');
        if (token) {
            config.headers.Authorization = 'Bearer ' + token;
        }
        return config;
    },
    (error) => {
        // 处理请求错误
        return Promise.reject(error);
    },
);

// 添加响应拦截器
axiosInstance.interceptors.response.use(
    // 请求成功
    (response: AxiosResponse) => {
        if (response.status !== 200) {
            errorHandle(response);
            return response;
        }

        return response;
    },
    // 请求失败
    (error: any) => {
        const {response} = error;
        if (response) {
            // 请求已发出，但是不在2xx的范围
            errorHandle(response);
            return response;
        }

        // 处理断网的情况
        // eg:请求超时或断网时，更新state的network状态
        // network状态在app.vue中控制着一个全局的断网提示组件的显示隐藏
        // 关于断网组件中的刷新重新获取数据，会在断网组件中说明
        message.error('网络连接异常,请稍后再试!');
    },
);

/**
 * http握手错误
 * @param res 响应回调,根据不同响应进行不同操作
 */
function errorHandle(res: any) {
    // 状态码判断
    switch (res.status) {
        case 401:
            break;
        case 403:
            break;
        case 404:
            message.error('请求的资源不存在');
            break;
        default:
            message.error(res.data?.message);
    }
}

export default axiosInstance;