// 使用项目的 request (axios) 替代原生 fetch
// 注意：这个文件在 src/api/generated/ 下，会被 generate-ts.sh（goctl api ts）整体覆盖生成，
// 上面这行 axios 替换和下面的 webapi.options 方法都是历史遗留的手工补丁，重新生成后需要手动重新应用，
// 否则会打回 goctl 默认的 fetch 实现、并再次缺失 options 方法（导致 admin.api 里 xxxOptions() 探测接口报错）
import request from '@/utils/request';

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

export const webapi = {
    get<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        // 去掉 /api 前缀，因为 request 的 baseURL 已包含
        const cleanUrl = url.replace(/^\/api/, '');
        // GET 请求使用 params（路径参数已在 URL 中，req 只包含查询参数）
        return request.get<T>(cleanUrl, req ? {params: req} : {}) as Promise<T>;
    },
    delete<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        const cleanUrl = url.replace(/^\/api/, '');
        // DELETE 请求需要传递请求体（包含 id）
        return request.delete<T>(cleanUrl, {data: req}) as Promise<T>;
    },
    put<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        const cleanUrl = url.replace(/^\/api/, '');
        return request.put<T>(cleanUrl, req) as Promise<T>;
    },
    post<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        const cleanUrl = url.replace(/^\/api/, '');
        return request.post<T>(cleanUrl, req) as Promise<T>;
    },
    patch<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        const cleanUrl = url.replace(/^\/api/, '');
        return request.patch<T>(cleanUrl, req) as Promise<T>;
    },
    options<T>(url: string, req?: unknown, config?: unknown): Promise<T> {
        const cleanUrl = url.replace(/^\/api/, '');
        return request.options<T>(cleanUrl, req ? {params: req} : {}) as Promise<T>;
    }
};

export default webapi
