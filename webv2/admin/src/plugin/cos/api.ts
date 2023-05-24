import lbstore from "@/plugin/api/lbstore";
import {message} from "ant-design-vue";

// 请求用到的参数
const Bucket = 'baifile-1309918034';
const Region = 'ap-guangzhou';
const protocol = location.protocol === 'https:' ? 'https:' : 'http:';
const prefix = protocol + '//' + Bucket + '.cos.' + Region + '.myqcloud.com/';  // prefix 用于拼接请求 u

export async function uploadFile(file: any): Promise<string> {
    console.log(file)
    let url = ""
    try {
        // 获取签名
        const key = 'dir/' + file?.file?.name; // 这里指定上传目录和文件名
        const resp = await lbstore.getSignature({
            method: 'PUT',
            name: key,
        })
        const SecurityToken = resp.session_token;
        const auth = String(resp.signature)
        url = prefix + camSafeUrlEncode(key).replace(/%2F/g, '/');
        const xhr = new XMLHttpRequest();
        xhr.open('PUT', url, true);
        xhr.setRequestHeader('Authorization', auth);
        SecurityToken && xhr.setRequestHeader('x-cos-security-token', SecurityToken);
        xhr.upload.onprogress = function (e) {
            console.log('上传进度 ' + (Math.round(e.loaded / e.total * 10000) / 100) + '%');
        };
        xhr.onload = function () {
            if (/^2\d\d$/.test('' + xhr.status)) {
                const ETag = xhr.getResponseHeader('etag');
                const rResp = lbstore.reportUploadFile({
                    file: {
                        id: undefined,
                        created_at: undefined,
                        updated_at: undefined,
                        deleted_at: undefined,
                        creator_uid: undefined,
                        file_name: file.file.name,
                        file_ext: undefined,
                        object_key: undefined,
                        sign_url: undefined,
                        url: url,
                        file_type: file.file.type,
                        size: file.file.size,
                    },
                })
                console.log(rResp,ETag)
            } else {
                console.log('文件 ' + key + ' 上传失败，状态码：' + xhr.status);
            }
        };
        xhr.onerror = function () {
            console.log('文件 ' + key + ' 上传失败，请检查是否没配置 CORS 跨域规则');
        };
        xhr.send(file?.file);
    } catch (error: any) {
        message.error(error)
    }
    return url
}

// 对更多字符编码的 url encode 格式
const camSafeUrlEncode = function (str: string) {
    return encodeURIComponent(str)
        .replace(/!/g, '%21')
        .replace(/'/g, '%27')
        .replace(/\(/g, '%28')
        .replace(/\)/g, '%29')
        .replace(/\*/g, '%2A');
};