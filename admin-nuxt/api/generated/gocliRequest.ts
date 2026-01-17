export type Method =
    | 'get'
    | 'GET'
    | 'delete'
    | 'DELETE'
    | 'head'
    | 'HEAD'
    | 'options'
    | 'OPTIONS'
    | 'post'
    | 'POST'
    | 'put'
    | 'PUT'
    | 'patch'
    | 'PATCH';

/**
 * Parse route parameters for responseType
 */
const reg = /:[a-z|A-Z]+/g;

export function parseParams(url: string): Array<string> {
    const ps = url.match(reg);
    if (!ps) {
        return [];
    }
    return ps.map((k) => k.replace(/:/, ''));
}

/**
 * Generate url and parameters
 * @param url
 * @param params
 */
export function genUrl(url: string, params: any) {
    if (!params) {
        return url;
    }

    const ps = parseParams(url);
    ps.forEach((k) => {
        const reg = new RegExp(`:${k}`);
        url = url.replace(reg, params[k]);
    });

    const path: Array<string> = [];
    for (const key of Object.keys(params)) {
        if (!ps.find((k) => k === key)) {
            path.push(`${key}=${params[key]}`);
        }
    }

    return url + (path.length > 0 ? `?${path.join('&')}` : '');
}

/**
 * 获取 baseURL（符合 Nuxt 3 规范）
 * 使用环境变量配置，统一在客户端和服务端使用完整 URL
 * 注意：useRuntimeConfig 只能在 setup 或 composable 中调用
 */
function getBaseURL(): string {
    try {
        const runtimeConfig = useRuntimeConfig()
        return runtimeConfig.public.apiBase || 'http://localhost:20000'
    } catch {
        // 如果不在 setup 上下文（例如在模块顶层），使用环境变量或默认值
        // 服务端可以直接访问 process.env
        if (typeof window === 'undefined') {
            return process.env.NUXT_PUBLIC_API_BASE || process.env.API_BASE_URL || 'http://localhost:20000'
        }
        // 客户端回退到默认值
        return 'http://localhost:20000'
    }
}

/**
 * 使用 Nuxt 3 的 $fetch 进行请求
 * 符合 Nuxt 3 开发规范
 */
export async function request<T = unknown>({
    method,
    url,
    data,
    config = {}
}: {
    method: Method;
    url: string;
    data?: unknown;
    config?: unknown;
}): Promise<T> {
    // 在 Nuxt 3 中，$fetch 会自动可用（自动导入）
    const baseURL = getBaseURL()
    // 如果 baseURL 为空字符串（使用代理），直接使用 url（相对路径）
    // 如果 baseURL 不为空，拼接完整 URL
    const fullUrl = url.startsWith('http') 
      ? url 
      : (baseURL ? `${baseURL}${url}` : url)
    
    const methodUpper = method.toLocaleUpperCase()
    const isGetOrDelete = methodUpper === 'GET' || methodUpper === 'DELETE' || methodUpper === 'HEAD' || methodUpper === 'OPTIONS'
    
    // GET/HEAD/DELETE/OPTIONS 请求不能有 body，数据通过 query 参数传递
    const fetchOptions: any = {
        method: methodUpper,
        credentials: 'include',
        // @ts-ignore
        ...config
    }
    
    if (isGetOrDelete) {
        // GET/DELETE 请求：如果有 data，应该已经通过 genUrl 处理为 query 参数
        // 不需要设置 body 和 Content-Type
    } else {
        // POST/PUT/PATCH 请求：设置 body 和 Content-Type
        fetchOptions.headers = {
            'Content-Type': 'application/json',
            // @ts-ignore
            ...(config?.headers || {})
        }
        if (data) {
            fetchOptions.body = JSON.stringify(data)
        }
    }
    
    return await $fetch<T>(fullUrl, fetchOptions)
}

function api<T>(
    method: Method = 'get',
    url: string,
    req: any,
    config?: unknown
): Promise<T> {
    const methodLower = method.toLocaleLowerCase() as Method
    const isGetOrDelete = methodLower === 'get' || methodLower === 'delete' || methodLower === 'head' || methodLower === 'options'
    
    // GET/DELETE/HEAD/OPTIONS 请求：将参数转换为 query string
    if (url.match(/:/) || isGetOrDelete) {
        // 处理路径参数和 query 参数
        const params = req?.params || req?.forms || req || {}
        url = genUrl(url, params)
    }
    
    // GET/DELETE/HEAD/OPTIONS 请求不需要 body，data 设为 undefined
    // POST/PUT/PATCH 请求需要 body
    const requestData = isGetOrDelete ? undefined : req

    switch (methodLower) {
        case 'get':
            return request<T>({method: 'get', url, data: requestData, config})
        case 'delete':
            return request<T>({method: 'delete', url, data: requestData, config})
        case 'put':
            return request<T>({method: 'put', url, data: requestData, config})
        case 'post':
            return request<T>({method: 'post', url, data: requestData, config})
        case 'patch':
            return request<T>({method: 'patch', url, data: requestData, config})
        case 'options':
            return request<T>({method: 'options', url, data: requestData, config})
        default:
            return request<T>({method: 'post', url, data: requestData, config})
    }
}

export const webapi = {
    get<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('get', url, req || {}, config);
    },
    delete<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('delete', url, req || {}, config);
    },
    put<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('put', url, req || {}, config);
    },
    post<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('post', url, req || {}, config);
    },
    patch<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('patch', url, req || {}, config);
    },
    options<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        return api<T>('options', url, req || {}, config);
    }
};

export default webapi
